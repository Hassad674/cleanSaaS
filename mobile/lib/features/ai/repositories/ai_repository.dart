import 'package:cleansaas_mobile/core/api/api_client.dart';
import 'package:cleansaas_mobile/features/ai/models/conversation.dart';
import 'package:cleansaas_mobile/features/ai/models/message.dart';

/// Repository handling all AI-related API communication.
///
/// Provides methods for CRUD operations on conversations and messages.
/// All network calls go through the centralized [ApiClient].
class AiRepository {
  final ApiClient _apiClient;

  AiRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Fetches all conversations for the current user.
  ///
  /// Supports pagination via [page] and [limit] parameters.
  Future<List<Conversation>> getConversations({
    int page = 1,
    int limit = 20,
  }) async {
    final response = await _apiClient.get(
      '/ai/conversations',
      queryParameters: {
        'page': page.toString(),
        'limit': limit.toString(),
      },
    );

    final List<dynamic> data = response.data['data'] as List<dynamic>;
    return data
        .map((json) => Conversation.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Creates a new conversation with an optional initial [title].
  Future<Conversation> createConversation({String? title}) async {
    final response = await _apiClient.post(
      '/ai/conversations',
      data: {
        if (title != null) 'title': title,
      },
    );

    return Conversation.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Deletes a conversation by [conversationId].
  Future<void> deleteConversation(String conversationId) async {
    await _apiClient.delete('/ai/conversations/$conversationId');
  }

  /// Fetches all messages for a given [conversationId].
  ///
  /// Messages are returned in chronological order.
  Future<List<Message>> getMessages(
    String conversationId, {
    int page = 1,
    int limit = 50,
  }) async {
    final response = await _apiClient.get(
      '/ai/conversations/$conversationId/messages',
      queryParameters: {
        'page': page.toString(),
        'limit': limit.toString(),
      },
    );

    final List<dynamic> data = response.data['data'] as List<dynamic>;
    return data
        .map((json) => Message.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Sends a user [content] message to a conversation and returns the
  /// AI assistant's response message.
  Future<Message> sendMessage(String conversationId, String content) async {
    final response = await _apiClient.post(
      '/ai/conversations/$conversationId/messages',
      data: {'content': content},
    );

    return Message.fromJson(response.data['data'] as Map<String, dynamic>);
  }
}
