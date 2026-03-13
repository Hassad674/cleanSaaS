/// Represents a single message within an AI conversation.
///
/// Messages have a [role] of either 'user' or 'assistant' to distinguish
/// who sent the message. Maps to the backend `messages` table.
class Message {
  final String id;
  final String conversationId;
  final String role;
  final String content;
  final DateTime createdAt;

  /// Role constant for user-sent messages.
  static const String roleUser = 'user';

  /// Role constant for AI-generated responses.
  static const String roleAssistant = 'assistant';

  const Message({
    required this.id,
    required this.conversationId,
    required this.role,
    required this.content,
    required this.createdAt,
  });

  bool get isUser => role == roleUser;
  bool get isAssistant => role == roleAssistant;

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: json['id'] as String,
      conversationId: json['conversation_id'] as String,
      role: json['role'] as String,
      content: json['content'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'conversation_id': conversationId,
      'role': role,
      'content': content,
      'created_at': createdAt.toIso8601String(),
    };
  }

  Message copyWith({
    String? id,
    String? conversationId,
    String? role,
    String? content,
    DateTime? createdAt,
  }) {
    return Message(
      id: id ?? this.id,
      conversationId: conversationId ?? this.conversationId,
      role: role ?? this.role,
      content: content ?? this.content,
      createdAt: createdAt ?? this.createdAt,
    );
  }
}
