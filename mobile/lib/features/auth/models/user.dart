/// User model matching the backend's [response.UserResponse].
///
/// Fields: id, email, name, role, avatar_url, email_verified.
/// The backend also returns created_at on some endpoints — we include it
/// as optional for forward-compatibility.
class User {
  const User({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
    this.avatarUrl = '',
    this.emailVerified = false,
    this.createdAt,
  });

  final String id;
  final String email;
  final String name;
  final String role;
  final String avatarUrl;
  final bool emailVerified;
  final DateTime? createdAt;

  /// Deserializes from the backend JSON shape:
  /// ```json
  /// {
  ///   "id": "uuid",
  ///   "email": "user@example.com",
  ///   "name": "John",
  ///   "role": "member",
  ///   "avatar_url": "",
  ///   "email_verified": false
  /// }
  /// ```
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      email: json['email'] as String,
      name: json['name'] as String,
      role: json['role'] as String? ?? 'member',
      avatarUrl: json['avatar_url'] as String? ?? '',
      emailVerified: json['email_verified'] as bool? ?? false,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'] as String)
          : null,
    );
  }

  /// Serializes to JSON for local caching.
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'role': role,
      'avatar_url': avatarUrl,
      'email_verified': emailVerified,
      if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
    };
  }

  /// Returns a copy with updated fields.
  User copyWith({
    String? id,
    String? email,
    String? name,
    String? role,
    String? avatarUrl,
    bool? emailVerified,
    DateTime? createdAt,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      name: name ?? this.name,
      role: role ?? this.role,
      avatarUrl: avatarUrl ?? this.avatarUrl,
      emailVerified: emailVerified ?? this.emailVerified,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  bool get isAdmin => role == 'admin';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is User && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() => 'User(id: $id, email: $email, name: $name)';
}
