import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:share_plus/share_plus.dart';
import 'package:cleansaas_mobile/features/blog/providers/blog_provider.dart';
import 'package:cleansaas_mobile/features/blog/models/blog_post.dart';

/// Full article view screen for a single blog post.
///
/// Displays cover image, title, author, date, full content, and a share button.
/// Fetches the post by slug using the [blogPostBySlugProvider].
class BlogPostScreen extends ConsumerWidget {
  final String slug;

  const BlogPostScreen({
    super.key,
    required this.slug,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final postAsync = ref.watch(blogPostBySlugProvider(slug));

    return Scaffold(
      body: postAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => _buildErrorState(context, ref, theme, error),
        data: (post) => _buildContent(context, theme, post),
      ),
    );
  }

  Widget _buildContent(
    BuildContext context,
    ThemeData theme,
    BlogPost post,
  ) {
    return CustomScrollView(
      slivers: [
        // Collapsing app bar with cover image
        SliverAppBar(
          expandedHeight: post.coverImageUrl != null ? 250 : 0,
          pinned: true,
          flexibleSpace: post.coverImageUrl != null
              ? FlexibleSpaceBar(
                  background: Image.network(
                    post.coverImageUrl!,
                    fit: BoxFit.cover,
                    errorBuilder: (_, __, ___) => Container(
                      color:
                          theme.colorScheme.primaryContainer.withOpacity(0.3),
                    ),
                  ),
                )
              : null,
          actions: [
            IconButton(
              icon: const Icon(Icons.share_outlined),
              onPressed: () => _sharePost(context, post),
              tooltip: 'Share',
            ),
          ],
        ),

        // Article content
        SliverToBoxAdapter(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Tags
                if (post.tags.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 12),
                    child: Wrap(
                      spacing: 8,
                      runSpacing: 6,
                      children: post.tags.map((tag) {
                        return Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 4,
                          ),
                          decoration: BoxDecoration(
                            color: theme.colorScheme.primary.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            tag,
                            style: theme.textTheme.labelSmall?.copyWith(
                              color: theme.colorScheme.primary,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  ),

                // Title
                Text(
                  post.title,
                  style: theme.textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    height: 1.3,
                  ),
                ),
                const SizedBox(height: 16),

                // Author row
                Row(
                  children: [
                    CircleAvatar(
                      radius: 18,
                      backgroundColor: theme.colorScheme.primaryContainer,
                      child: Text(
                        post.authorName.isNotEmpty
                            ? post.authorName[0].toUpperCase()
                            : '?',
                        style: theme.textTheme.titleSmall?.copyWith(
                          color: theme.colorScheme.onPrimaryContainer,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          post.authorName,
                          style: theme.textTheme.bodyMedium?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        Text(
                          '${post.formattedDate}  ·  ${post.readTime}',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color:
                                theme.colorScheme.onSurface.withOpacity(0.5),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
                const SizedBox(height: 24),

                // Divider
                Divider(color: theme.colorScheme.outlineVariant),
                const SizedBox(height: 24),

                // Article body
                // Renders content as plain text paragraphs. For HTML content,
                // consider using a package like flutter_html or flutter_widget_from_html.
                _buildArticleBody(theme, post.content),
                const SizedBox(height: 40),

                // Share button at bottom
                Center(
                  child: OutlinedButton.icon(
                    onPressed: () => _sharePost(context, post),
                    icon: const Icon(Icons.share_outlined),
                    label: const Text('Share this article'),
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 24,
                        vertical: 12,
                      ),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 32),
              ],
            ),
          ),
        ),
      ],
    );
  }

  /// Renders the article body as formatted text paragraphs.
  ///
  /// Splits content by double newlines into paragraphs. For rich HTML
  /// content rendering, integrate a package like `flutter_html`.
  Widget _buildArticleBody(ThemeData theme, String content) {
    final paragraphs = content.split(RegExp(r'\n\n+'));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: paragraphs.map((paragraph) {
        final trimmed = paragraph.trim();
        if (trimmed.isEmpty) return const SizedBox.shrink();

        // Detect headings (lines starting with #).
        if (trimmed.startsWith('# ')) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 12, top: 8),
            child: Text(
              trimmed.substring(2),
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          );
        }
        if (trimmed.startsWith('## ')) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 10, top: 8),
            child: Text(
              trimmed.substring(3),
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          );
        }
        if (trimmed.startsWith('### ')) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 8, top: 6),
            child: Text(
              trimmed.substring(4),
              style: theme.textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          );
        }

        // Regular paragraph.
        return Padding(
          padding: const EdgeInsets.only(bottom: 16),
          child: Text(
            trimmed,
            style: theme.textTheme.bodyLarge?.copyWith(
              height: 1.7,
              color: theme.colorScheme.onSurface.withOpacity(0.85),
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildErrorState(
    BuildContext context,
    WidgetRef ref,
    ThemeData theme,
    Object error,
  ) {
    return Scaffold(
      appBar: AppBar(),
      body: Center(
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
                'Failed to load article',
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
                onPressed: () =>
                    ref.invalidate(blogPostBySlugProvider(slug)),
                icon: const Icon(Icons.refresh),
                label: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _sharePost(BuildContext context, BlogPost post) {
    // Share the post title and a deep link or web URL.
    Share.share(
      '${post.title}\n\nRead more: https://cleansaas.com/blog/${post.slug}',
      subject: post.title,
    );
  }
}
