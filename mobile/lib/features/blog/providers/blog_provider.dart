import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/auth/providers/auth_provider.dart';
import 'package:cleansaas_mobile/features/blog/models/blog_post.dart';
import 'package:cleansaas_mobile/features/blog/repositories/blog_repository.dart';

/// Provider for the [BlogRepository] instance.
final blogRepositoryProvider = Provider<BlogRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return BlogRepository(apiClient);
});

/// Provider tracking the currently selected tag filter.
///
/// null means "all posts" (no filter).
final selectedTagProvider = StateProvider<String?>((ref) {
  return null;
});

/// Async provider that fetches available blog tags.
final blogTagsProvider = FutureProvider<List<String>>((ref) async {
  final repository = ref.watch(blogRepositoryProvider);
  return repository.getTags();
});

/// Async notifier managing the blog posts list with tag filtering and pagination.
final blogPostsProvider =
    AsyncNotifierProvider<BlogPostsNotifier, List<BlogPost>>(
  BlogPostsNotifier.new,
);

/// Notifier managing blog posts with tag filtering support.
///
/// Automatically refetches when the selected tag changes.
class BlogPostsNotifier extends AsyncNotifier<List<BlogPost>> {
  int _currentPage = 1;
  bool _hasMore = true;
  static const _pageSize = 20;

  @override
  Future<List<BlogPost>> build() async {
    final repository = ref.watch(blogRepositoryProvider);
    final tag = ref.watch(selectedTagProvider);

    _currentPage = 1;
    _hasMore = true;

    final posts = await repository.getPosts(
      tag: tag,
      page: 1,
      limit: _pageSize,
    );

    _hasMore = posts.length >= _pageSize;
    return posts;
  }

  /// Loads the next page of posts. No-op if all posts have been loaded.
  Future<void> loadMore() async {
    if (!_hasMore) return;

    final repository = ref.read(blogRepositoryProvider);
    final tag = ref.read(selectedTagProvider);
    final currentPosts = state.valueOrNull ?? [];

    _currentPage++;

    try {
      final newPosts = await repository.getPosts(
        tag: tag,
        page: _currentPage,
        limit: _pageSize,
      );

      _hasMore = newPosts.length >= _pageSize;
      state = AsyncValue.data([...currentPosts, ...newPosts]);
    } catch (e, stack) {
      // Revert page increment on failure.
      _currentPage--;
      state = AsyncValue.error(e, stack);
    }
  }

  /// Whether there are more posts to load.
  bool get hasMore => _hasMore;

  /// Forces a full refresh of the post list.
  Future<void> refresh() async {
    ref.invalidateSelf();
  }
}

/// Family provider to fetch a single blog post by slug.
///
/// Caches each post independently to avoid refetching when navigating back.
final blogPostBySlugProvider =
    FutureProvider.family<BlogPost, String>((ref, slug) async {
  final repository = ref.watch(blogRepositoryProvider);
  return repository.getPostBySlug(slug);
});
