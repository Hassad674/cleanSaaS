import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/auth/providers/auth_provider.dart';
import 'package:cleansaas_mobile/features/settings/repositories/user_repository.dart';

/// Provider for the [UserRepository] instance.
final userRepositoryProvider = Provider<UserRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return UserRepository(apiClient);
});

/// Async provider that fetches and caches the current user's profile.
///
/// Invalidate this provider to trigger a fresh fetch from the API.
final userProfileProvider = AsyncNotifierProvider<UserProfileNotifier, UserProfile>(
  UserProfileNotifier.new,
);

/// Notifier managing the user profile state.
///
/// Handles fetching, updating profile fields, changing password,
/// uploading avatar, and deleting the account.
class UserProfileNotifier extends AsyncNotifier<UserProfile> {
  @override
  Future<UserProfile> build() async {
    final repository = ref.watch(userRepositoryProvider);
    return repository.getProfile();
  }

  /// Updates the user's name and/or email.
  Future<void> updateProfile({String? name, String? email}) async {
    final repository = ref.read(userRepositoryProvider);
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      return repository.updateProfile(name: name, email: email);
    });
  }

  /// Changes the user's password.
  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    final repository = ref.read(userRepositoryProvider);
    await repository.changePassword(
      currentPassword: currentPassword,
      newPassword: newPassword,
    );
  }

  /// Uploads a new avatar image.
  Future<void> uploadAvatar({
    required List<int> imageBytes,
    required String filename,
  }) async {
    final repository = ref.read(userRepositoryProvider);
    final avatarUrl = await repository.uploadAvatar(
      imageBytes: imageBytes,
      filename: filename,
    );

    // Update the local state with the new avatar URL.
    final current = state.valueOrNull;
    if (current != null) {
      state = AsyncValue.data(current.copyWith(avatarUrl: avatarUrl));
    }
  }

  /// Deletes the user's account permanently.
  Future<void> deleteAccount({required String password}) async {
    final repository = ref.read(userRepositoryProvider);
    await repository.deleteAccount(password: password);
  }

  /// Forces a refresh of the profile from the API.
  Future<void> refresh() async {
    ref.invalidateSelf();
  }
}

/// Provider for the dark mode toggle state.
///
/// Persists the preference locally. Defaults to system theme.
final darkModeProvider = StateProvider<bool?>((ref) {
  // null = follow system, true = dark, false = light
  return null;
});

/// Provider for the notifications toggle state.
final notificationsEnabledProvider = StateProvider<bool>((ref) {
  return true;
});
