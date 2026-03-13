import '../../../core/api/api_client.dart';
import '../../../core/api/api_exceptions.dart';
import '../models/user.dart';

/// Holds the token + user pair returned by login and register endpoints.
class AuthResult {
  const AuthResult({required this.token, required this.user});

  final String token;
  final User user;
}

/// Handles all HTTP calls to the backend auth endpoints.
///
/// Endpoints (matching [router.go]):
/// - POST /auth/register
/// - POST /auth/login
/// - POST /auth/forgot-password
/// - POST /auth/reset-password
/// - POST /auth/verify-email
/// - POST /auth/resend-verification (authenticated)
class AuthRepository {
  const AuthRepository({required ApiClient apiClient}) : _api = apiClient;

  final ApiClient _api;

  /// Registers a new user.
  ///
  /// Returns [AuthResult] with the JWT token and user data on success.
  /// Throws [ConflictException] if the email is already taken.
  /// Throws [BadRequestException] if fields are missing or invalid.
  Future<AuthResult> register({
    required String email,
    required String name,
    required String password,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      '/auth/register',
      data: {
        'email': email,
        'name': name,
        'password': password,
      },
    );

    return _parseAuthResponse(response.data!);
  }

  /// Authenticates an existing user.
  ///
  /// Returns [AuthResult] on success.
  /// Throws [UnauthorizedException] if credentials are wrong.
  Future<AuthResult> login({
    required String email,
    required String password,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      '/auth/login',
      data: {
        'email': email,
        'password': password,
      },
    );

    return _parseAuthResponse(response.data!);
  }

  /// Sends a password-reset email if the account exists.
  ///
  /// The backend always returns 200 to avoid leaking user existence.
  Future<void> forgotPassword({required String email}) async {
    await _api.post<Map<String, dynamic>>(
      '/auth/forgot-password',
      data: {'email': email},
    );
  }

  /// Resets the password using a token received via email.
  ///
  /// Throws [BadRequestException] if the token is invalid or expired.
  Future<void> resetPassword({
    required String token,
    required String password,
  }) async {
    await _api.post<Map<String, dynamic>>(
      '/auth/reset-password',
      data: {
        'token': token,
        'password': password,
      },
    );
  }

  /// Verifies a user's email address using the token sent via email.
  Future<void> verifyEmail({required String token}) async {
    await _api.post<Map<String, dynamic>>(
      '/auth/verify-email',
      data: {'token': token},
    );
  }

  /// Resends the verification email for the currently authenticated user.
  Future<void> resendVerification() async {
    await _api.post<Map<String, dynamic>>('/auth/resend-verification');
  }

  /// Fetches the currently authenticated user's profile.
  ///
  /// Used on app launch to verify the stored token is still valid.
  Future<User> getProfile() async {
    final response = await _api.get<Map<String, dynamic>>('/users/me');
    return User.fromJson(response.data!);
  }

  // ---------------------------------------------------------------------------
  // Helpers
  // ---------------------------------------------------------------------------

  /// Parses the backend's `{ "token": "...", "user": {...} }` response.
  AuthResult _parseAuthResponse(Map<String, dynamic> data) {
    return AuthResult(
      token: data['token'] as String,
      user: User.fromJson(data['user'] as Map<String, dynamic>),
    );
  }
}
