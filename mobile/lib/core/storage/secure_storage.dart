import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../../config/constants.dart';

/// Abstraction over [FlutterSecureStorage] for persisting sensitive data
/// such as JWT tokens and cached user JSON.
class SecureStorageService {
  SecureStorageService({FlutterSecureStorage? storage})
      : _storage = storage ??
            const FlutterSecureStorage(
              aOptions: AndroidOptions(encryptedSharedPreferences: true),
              iOptions: IOSOptions(
                accessibility: KeychainAccessibility.first_unlock,
              ),
            );

  final FlutterSecureStorage _storage;

  // ---------------------------------------------------------------------------
  // Token
  // ---------------------------------------------------------------------------

  /// Persists the JWT auth token.
  Future<void> saveToken(String token) async {
    await _storage.write(key: AppConstants.tokenKey, value: token);
  }

  /// Returns the stored JWT token, or `null` if none exists.
  Future<String?> getToken() async {
    return _storage.read(key: AppConstants.tokenKey);
  }

  /// Removes the stored JWT token.
  Future<void> deleteToken() async {
    await _storage.delete(key: AppConstants.tokenKey);
  }

  // ---------------------------------------------------------------------------
  // User JSON cache
  // ---------------------------------------------------------------------------

  /// Persists a JSON-serializable user map for offline access.
  Future<void> saveUser(Map<String, dynamic> userJson) async {
    await _storage.write(
      key: AppConstants.userKey,
      value: jsonEncode(userJson),
    );
  }

  /// Returns the cached user JSON map, or `null` if nothing is stored.
  Future<Map<String, dynamic>?> getUser() async {
    final raw = await _storage.read(key: AppConstants.userKey);
    if (raw == null) return null;
    return jsonDecode(raw) as Map<String, dynamic>;
  }

  /// Removes the cached user JSON.
  Future<void> deleteUser() async {
    await _storage.delete(key: AppConstants.userKey);
  }

  // ---------------------------------------------------------------------------
  // Utilities
  // ---------------------------------------------------------------------------

  /// Returns `true` if a valid token is stored.
  Future<bool> hasToken() async {
    final token = await getToken();
    return token != null && token.isNotEmpty;
  }

  /// Clears all stored credentials (token + user cache).
  Future<void> clearAll() async {
    await Future.wait([
      deleteToken(),
      deleteUser(),
    ]);
  }
}
