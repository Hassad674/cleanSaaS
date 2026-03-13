import 'package:dio/dio.dart';
import 'package:cleansaas_mobile/core/api/api_client.dart';

/// Data model for the authenticated user's profile.
class UserProfile {
  final String id;
  final String name;
  final String email;
  final String? avatarUrl;
  final DateTime createdAt;
  final DateTime updatedAt;

  const UserProfile({
    required this.id,
    required this.name,
    required this.email,
    this.avatarUrl,
    required this.createdAt,
    required this.updatedAt,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      id: json['id'] as String,
      name: json['name'] as String,
      email: json['email'] as String,
      avatarUrl: json['avatar_url'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'email': email,
      'avatar_url': avatarUrl,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  UserProfile copyWith({
    String? name,
    String? email,
    String? avatarUrl,
  }) {
    return UserProfile(
      id: id,
      name: name ?? this.name,
      email: email ?? this.email,
      avatarUrl: avatarUrl ?? this.avatarUrl,
      createdAt: createdAt,
      updatedAt: DateTime.now(),
    );
  }
}

/// Repository handling all user profile API operations.
///
/// Provides methods to fetch, update, change password, upload avatar,
/// and delete the current user's account.
class UserRepository {
  final ApiClient _apiClient;

  UserRepository(this._apiClient);

  /// Fetches the current authenticated user's profile.
  Future<UserProfile> getProfile() async {
    final response = await _apiClient.get('/users/me');
    return UserProfile.fromJson(response.data as Map<String, dynamic>);
  }

  /// Updates the user's profile fields (name and/or email).
  Future<UserProfile> updateProfile({
    String? name,
    String? email,
  }) async {
    final body = <String, dynamic>{};
    if (name != null) body['name'] = name;
    if (email != null) body['email'] = email;

    final response = await _apiClient.patch('/users/me', data: body);
    return UserProfile.fromJson(response.data as Map<String, dynamic>);
  }

  /// Changes the user's password.
  ///
  /// Requires the current password for verification and the new password.
  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    await _apiClient.put('/users/me/password', data: {
      'current_password': currentPassword,
      'new_password': newPassword,
    });
  }

  /// Uploads a new avatar image for the user.
  ///
  /// [imageBytes] is the raw image data, [filename] is the original filename.
  /// Returns the URL of the uploaded avatar.
  Future<String> uploadAvatar({
    required List<int> imageBytes,
    required String filename,
  }) async {
    final formData = FormData.fromMap({
      'avatar': MultipartFile.fromBytes(imageBytes, filename: filename),
    });

    final response = await _apiClient.post(
      '/users/me/avatar',
      data: formData,
    );
    return response.data['avatar_url'] as String;
  }

  /// Permanently deletes the user's account.
  ///
  /// This action is irreversible. Requires password confirmation.
  Future<void> deleteAccount({required String password}) async {
    await _apiClient.delete('/users/me', data: {
      'password': password,
    });
  }
}
