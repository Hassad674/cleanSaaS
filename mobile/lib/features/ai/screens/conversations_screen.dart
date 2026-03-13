import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cleansaas_mobile/features/ai/providers/ai_provider.dart';
import 'package:cleansaas_mobile/features/ai/widgets/conversation_tile.dart';

/// Screen that lists all AI conversations for the current user.
///
/// Supports pull-to-refresh, infinite scroll pagination, swipe-to-delete,
/// and a floating action button to create new conversations.
class ConversationsScreen extends ConsumerStatefulWidget {
  const ConversationsScreen({super.key});

  @override
  ConsumerState<ConversationsScreen> createState() =>
      _ConversationsScreenState();
}

class _ConversationsScreenState extends ConsumerState<ConversationsScreen> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);

    // Load conversations on first build.
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(aiProvider).loadConversations();
    });
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent - 200) {
      ref.read(aiProvider).loadMoreConversations();
    }
  }

  Future<void> _onRefresh() async {
    await ref.read(aiProvider).loadConversations();
  }

  Future<void> _createConversation() async {
    final provider = ref.read(aiProvider);
    final conversation = await provider.createConversation();
    if (conversation != null && mounted) {
      context.push('/ai/${conversation.id}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final ai = ref.watch(aiProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'AI Chat',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _createConversation,
        child: const Icon(Icons.add),
      ),
      body: _buildBody(theme, ai),
    );
  }

  Widget _buildBody(ThemeData theme, AiProvider ai) {
    if (ai.isLoadingConversations && ai.conversations.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (ai.conversationsError != null && ai.conversations.isEmpty) {
      return _buildErrorState(theme, ai);
    }

    if (ai.conversations.isEmpty) {
      return _buildEmptyState(theme);
    }

    return RefreshIndicator(
      onRefresh: _onRefresh,
      child: ListView.builder(
        controller: _scrollController,
        padding: const EdgeInsets.symmetric(vertical: 8),
        itemCount: ai.conversations.length + (ai.hasMoreConversations ? 1 : 0),
        itemBuilder: (context, index) {
          if (index >= ai.conversations.length) {
            return const Padding(
              padding: EdgeInsets.all(16),
              child: Center(child: CircularProgressIndicator()),
            );
          }

          final conversation = ai.conversations[index];
          return ConversationTile(
            conversation: conversation,
            onTap: () {
              ai.selectConversation(conversation);
              context.push('/ai/${conversation.id}');
            },
            onDelete: () => ai.deleteConversation(conversation.id),
          );
        },
      ),
    );
  }

  Widget _buildErrorState(ThemeData theme, AiProvider ai) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.error_outline,
              size: 48,
              color: theme.colorScheme.error,
            ),
            const SizedBox(height: 16),
            Text(
              ai.conversationsError!,
              style: theme.textTheme.bodyLarge,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: () => ai.loadConversations(),
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(ThemeData theme) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.chat_bubble_outline,
              size: 64,
              color: theme.colorScheme.onSurface.withOpacity(0.3),
            ),
            const SizedBox(height: 16),
            Text(
              'No conversations yet',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Start a new conversation with the AI assistant',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.5),
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: _createConversation,
              icon: const Icon(Icons.add),
              label: const Text('New Conversation'),
            ),
          ],
        ),
      ),
    );
  }
}
