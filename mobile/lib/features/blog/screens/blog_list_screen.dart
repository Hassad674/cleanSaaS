import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/blog/providers/blog_provider.dart';
import 'package:cleansaas_mobile/features/blog/widgets/post_card.dart';
import 'package:cleansaas_mobile/features/blog/widgets/tag_chip.dart';
import 'package:cleansaas_mobile/features/blog/screens/blog_post_screen.dart';

/// Scrollable list of blog post cards with tag filter chips at the top.
///
/// Supports pull-to-refresh and infinite scroll pagination.
class BlogListScreen extends ConsumerStatefulWidget {
  const BlogListScreen({super.key});

  @override
  ConsumerState<BlogListScreen> createState() => _BlogListScreenState();
}

class _BlogListScreenState extends ConsumerState<BlogListScreen> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
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
      ref.read(blogPostsProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final postsAsync = ref.watch(blogPostsProvider);
    final tagsAsync = ref.watch(blogTagsProvider);
    final selectedTag = ref.watch(selectedTagProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Blog',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(blogPostsProvider);
          ref.invalidate(blogTagsProvider);
        },
        child: CustomScrollView(
          controller: _scrollController,
          slivers: [
            // Tag filter chips
            SliverToBoxAdapter(
              child: tagsAsync.when(
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
                data: (tags) {
                  if (tags.isEmpty) return const SizedBox.shrink();
                  return _buildTagFilters(theme, tags, selectedTag);
                },
              ),
            ),

            // Posts list
            postsAsync.when(
              loading: () => const SliverFillRemaining(
                child: Center(child: CircularProgressIndicator()),
              ),
              error: (error, _) => SliverFillRemaining(
                child: _buildErrorState(theme, error),
              ),
              data: (posts) {
                if (posts.isEmpty) {
                  return SliverFillRemaining(
                    child: _buildEmptyState(theme),
                  );
                }

                return SliverPadding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  sliver: SliverList.builder(
                    itemCount: posts.length + 1, // +1 for loading indicator
                    itemBuilder: (context, index) {
                      if (index == posts.length) {
                        return _buildLoadMoreIndicator(ref);
                      }

                      final post = posts[index];
                      return Padding(
                        padding: const EdgeInsets.only(bottom: 16),
                        child: PostCard(
                          post: post,
                          onTap: () => _navigateToPost(context, post.slug),
                        ),
                      );
                    },
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTagFilters(
    ThemeData theme,
    List<String> tags,
    String? selectedTag,
  ) {
    return SizedBox(
      height: 52,
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        itemCount: tags.length + 1, // +1 for "All" chip
        itemBuilder: (context, index) {
          if (index == 0) {
            return Padding(
              padding: const EdgeInsets.only(right: 8),
              child: TagChip(
                tag: 'All',
                isSelected: selectedTag == null,
                onTap: () {
                  ref.read(selectedTagProvider.notifier).state = null;
                },
              ),
            );
          }

          final tag = tags[index - 1];
          return Padding(
            padding: const EdgeInsets.only(right: 8),
            child: TagChip(
              tag: tag,
              isSelected: selectedTag == tag,
              onTap: () {
                ref.read(selectedTagProvider.notifier).state =
                    selectedTag == tag ? null : tag;
              },
            ),
          );
        },
      ),
    );
  }

  Widget _buildEmptyState(ThemeData theme) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.article_outlined,
              size: 48,
              color: theme.colorScheme.onSurface.withOpacity(0.3),
            ),
            const SizedBox(height: 16),
            Text(
              'No posts found',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'Check back later for new content.',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.6),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildErrorState(ThemeData theme, Object error) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.error_outline,
              size: 48,
              color: theme.colorScheme.error,
            ),
            const SizedBox(height: 16),
            Text(
              'Failed to load posts',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              error.toString(),
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.6),
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: () => ref.invalidate(blogPostsProvider),
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLoadMoreIndicator(WidgetRef ref) {
    final notifier = ref.read(blogPostsProvider.notifier);
    if (!notifier.hasMore) {
      return const SizedBox(height: 32);
    }

    return const Padding(
      padding: EdgeInsets.symmetric(vertical: 24),
      child: Center(child: CircularProgressIndicator()),
    );
  }

  void _navigateToPost(BuildContext context, String slug) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => BlogPostScreen(slug: slug),
      ),
    );
  }
}
