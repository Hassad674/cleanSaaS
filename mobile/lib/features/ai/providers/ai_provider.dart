import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/ai/models/conversation.dart';
import 'package:cleansaas_mobile/features/ai/models/message.dart';
import 'package:cleansaas_mobile/features/ai/repositories/ai_repository.dart';

/// State management for the AI chat feature.
///
/// Manages the conversations list, current conversation, messages,
/// and sending state. Uses [ChangeNotifier] for Riverpod integration.
class AiProvider extends ChangeNotifier {
  final AiRepository _repository;

  AiProvider({required AiRepository repository}) : _repository = repository;

  // -- Conversations state --

  List<Conversation> _conversations = [];
  List<Conversation> get conversations => _conversations;

  bool _isLoadingConversations = false;
  bool get isLoadingConversations => _isLoadingConversations;

  String? _conversationsError;
  String? get conversationsError => _conversationsError;

  int _conversationsPage = 1;
  bool _hasMoreConversations = true;
  bool get hasMoreConversations => _hasMoreConversations;

  // -- Current conversation / messages state --

  Conversation? _currentConversation;
  Conversation? get currentConversation => _currentConversation;

  List<Message> _messages = [];
  List<Message> get messages => _messages;

  bool _isLoadingMessages = false;
  bool get isLoadingMessages => _isLoadingMessages;

  bool _isSending = false;
  bool get isSending => _isSending;

  String? _messagesError;
  String? get messagesError => _messagesError;

  // -- Conversations --

  /// Loads the first page of conversations, replacing any existing data.
  Future<void> loadConversations() async {
    _isLoadingConversations = true;
    _conversationsError = null;
    _conversationsPage = 1;
    notifyListeners();

    try {
      _conversations = await _repository.getConversations(page: 1);
      _hasMoreConversations = _conversations.length >= 20;
    } catch (e) {
      _conversationsError = 'Failed to load conversations. Please try again.';
    } finally {
      _isLoadingConversations = false;
      notifyListeners();
    }
  }

  /// Loads the next page of conversations and appends them.
  Future<void> loadMoreConversations() async {
    if (_isLoadingConversations || !_hasMoreConversations) return;

    _isLoadingConversations = true;
    notifyListeners();

    try {
      _conversationsPage++;
      final more = await _repository.getConversations(
        page: _conversationsPage,
      );
      _conversations = [..._conversations, ...more];
      _hasMoreConversations = more.length >= 20;
    } catch (e) {
      _conversationsPage--;
      _conversationsError = 'Failed to load more conversations.';
    } finally {
      _isLoadingConversations = false;
      notifyListeners();
    }
  }

  /// Creates a new conversation and navigates to it by setting it as current.
  Future<Conversation?> createConversation({String? title}) async {
    try {
      final conversation = await _repository.createConversation(title: title);
      _conversations = [conversation, ..._conversations];
      _currentConversation = conversation;
      _messages = [];
      notifyListeners();
      return conversation;
    } catch (e) {
      _conversationsError = 'Failed to create conversation.';
      notifyListeners();
      return null;
    }
  }

  /// Deletes a conversation by [id] and removes it from the local list.
  Future<void> deleteConversation(String id) async {
    try {
      await _repository.deleteConversation(id);
      _conversations = _conversations.where((c) => c.id != id).toList();
      if (_currentConversation?.id == id) {
        _currentConversation = null;
        _messages = [];
      }
      notifyListeners();
    } catch (e) {
      _conversationsError = 'Failed to delete conversation.';
      notifyListeners();
    }
  }

  // -- Messages --

  /// Selects a conversation and loads its messages.
  Future<void> selectConversation(Conversation conversation) async {
    _currentConversation = conversation;
    _messages = [];
    _messagesError = null;
    _isLoadingMessages = true;
    notifyListeners();

    try {
      _messages = await _repository.getMessages(conversation.id);
    } catch (e) {
      _messagesError = 'Failed to load messages.';
    } finally {
      _isLoadingMessages = false;
      notifyListeners();
    }
  }

  /// Sends a message in the current conversation.
  ///
  /// Adds the user message optimistically, then appends the AI response
  /// when it arrives. On failure, removes the optimistic message.
  Future<void> sendMessage(String content) async {
    if (_currentConversation == null || _isSending) return;

    _isSending = true;
    _messagesError = null;

    // Optimistic user message.
    final optimisticMessage = Message(
      id: 'temp-${DateTime.now().millisecondsSinceEpoch}',
      conversationId: _currentConversation!.id,
      role: Message.roleUser,
      content: content,
      createdAt: DateTime.now(),
    );
    _messages = [..._messages, optimisticMessage];
    notifyListeners();

    try {
      final response = await _repository.sendMessage(
        _currentConversation!.id,
        content,
      );

      // Replace optimistic message with server-confirmed version and add AI reply.
      _messages = [
        ..._messages.where((m) => m.id != optimisticMessage.id),
        // The API may return the user message + assistant message, or just the
        // assistant message. We add the assistant response here.
        optimisticMessage.copyWith(
          id: 'confirmed-${DateTime.now().millisecondsSinceEpoch}',
        ),
        response,
      ];
    } catch (e) {
      // Remove optimistic message on failure.
      _messages = _messages.where((m) => m.id != optimisticMessage.id).toList();
      _messagesError = 'Failed to send message. Please try again.';
    } finally {
      _isSending = false;
      notifyListeners();
    }
  }

  /// Clears the current conversation selection.
  void clearCurrentConversation() {
    _currentConversation = null;
    _messages = [];
    _messagesError = null;
    notifyListeners();
  }
}

/// Riverpod provider for [AiProvider].
///
/// Requires an [AiRepository] to be provided in the widget tree.
final aiProvider = ChangeNotifierProvider<AiProvider>((ref) {
  throw UnimplementedError(
    'aiProvider must be overridden with a valid AiRepository.',
  );
});
