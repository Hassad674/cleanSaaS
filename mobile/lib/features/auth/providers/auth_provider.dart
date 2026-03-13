import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api/api_client.dart';
import '../../../core/api/api_exceptions.dart';
import '../../../core/storage/secure_storage.dart';
import '../models/user.dart';
import '../repositories/auth_repository.dart';

// ---------------------------------------------------------------------------
// Singleton service providers
// ---------------------------------------------------------------------------

/// Provides the [SecureStorageService] singleton.
final secureStorageProvider = Provider<SecureStorageService>((ref) {
  return SecureStorageService();
});

/// Provides the [ApiClient] singleton, injecting secure storage for auth.
final apiClientProvider = Provider<ApiClient>((ref) {
  final storage = ref.watch(secureStorageProvider);
  return ApiClient(storage: storage);
});

/// Provides the [AuthRepository] singleton.
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return AuthRepository(apiClient: apiClient);
});

// ---------------------------------------------------------------------------
// Auth state
// ---------------------------------------------------------------------------

/// Possible authentication states.
enum AuthStatus { loading, authenticated, unauthenticated }

/// Immutable snapshot of the authentication state.
@immutable
class AuthState {
  const AuthState({
    required this.status,
    this.user,
    this.errorMessage,
    this.isSubmitting = false,
  });

  final AuthStatus status;
  final User? user;
  final String? errorMessage;

  /// True while a login/register/forgot-password request is in-flight.
  final bool isSubmitting;

  const AuthState.initial()
      : status = AuthStatus.loading,
        user = null,
        errorMessage = null,
        isSubmitting = false;

  AuthState copyWith({
    AuthStatus? status,
    User? user,
    String? errorMessage,
    bool? isSubmitting,
  }) {
    return AuthState(
      status: status ?? this.status,
      user: user ?? this.user,
      errorMessage: errorMessage,
      isSubmitting: isSubmitting ?? this.isSubmitting,
    );
  }
}

// ---------------------------------------------------------------------------
// Auth notifier
// ---------------------------------------------------------------------------

/// Manages authentication state: login, register, logout, and session restore.
class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier({
    required AuthRepository repository,
    required SecureStorageService storage,
  })  : _repository = repository,
        _storage = storage,
        super(const AuthState.initial()) {
    _tryRestoreSession();
  }

  final AuthRepository _repository;
  final SecureStorageService _storage;

  /// Attempts to restore a session from stored credentials on app start.
  Future<void> _tryRestoreSession() async {
    try {
      final hasToken = await _storage.hasToken();
      if (!hasToken) {
        state = state.copyWith(status: AuthStatus.unauthenticated);
        return;
      }

      // Verify the token is still valid by hitting /users/me.
      final user = await _repository.getProfile();
      await _storage.saveUser(user.toJson());

      state = AuthState(
        status: AuthStatus.authenticated,
        user: user,
      );
    } on ApiException {
      // Token invalid or expired — clear and go to login.
      await _storage.clearAll();
      state = state.copyWith(status: AuthStatus.unauthenticated);
    } catch (_) {
      // Network error — try cached user as fallback.
      final cachedUser = await _storage.getUser();
      if (cachedUser != null) {
        state = AuthState(
          status: AuthStatus.authenticated,
          user: User.fromJson(cachedUser),
        );
      } else {
        state = state.copyWith(status: AuthStatus.unauthenticated);
      }
    }
  }

  /// Logs in with email and password.
  Future<bool> login({
    required String email,
    required String password,
  }) async {
    state = state.copyWith(isSubmitting: true, errorMessage: null);

    try {
      final result = await _repository.login(
        email: email,
        password: password,
      );

      await _storage.saveToken(result.token);
      await _storage.saveUser(result.user.toJson());

      state = AuthState(
        status: AuthStatus.authenticated,
        user: result.user,
      );
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: e.message,
      );
      return false;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: 'An unexpected error occurred',
      );
      return false;
    }
  }

  /// Registers a new account.
  Future<bool> register({
    required String email,
    required String name,
    required String password,
  }) async {
    state = state.copyWith(isSubmitting: true, errorMessage: null);

    try {
      final result = await _repository.register(
        email: email,
        name: name,
        password: password,
      );

      await _storage.saveToken(result.token);
      await _storage.saveUser(result.user.toJson());

      state = AuthState(
        status: AuthStatus.authenticated,
        user: result.user,
      );
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: e.message,
      );
      return false;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: 'An unexpected error occurred',
      );
      return false;
    }
  }

  /// Sends a forgot-password request. Always succeeds from the user's
  /// perspective (backend does not leak user existence).
  Future<bool> forgotPassword({required String email}) async {
    state = state.copyWith(isSubmitting: true, errorMessage: null);

    try {
      await _repository.forgotPassword(email: email);
      state = state.copyWith(isSubmitting: false);
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: e.message,
      );
      return false;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: 'An unexpected error occurred',
      );
      return false;
    }
  }

  /// Logs out: clears stored credentials and resets to unauthenticated.
  Future<void> logout() async {
    await _storage.clearAll();
    state = const AuthState(status: AuthStatus.unauthenticated);
  }

  /// Clears any displayed error message.
  void clearError() {
    state = state.copyWith(errorMessage: null);
  }
}

// ---------------------------------------------------------------------------
// Provider
// ---------------------------------------------------------------------------

/// The main auth state provider.
///
/// Usage:
/// ```dart
/// final authState = ref.watch(authProvider);
/// final authNotifier = ref.read(authProvider.notifier);
/// ```
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  final storage = ref.watch(secureStorageProvider);
  return AuthNotifier(repository: repository, storage: storage);
});
