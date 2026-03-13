/// Represents a blog post with all its metadata and content.
class BlogPost {
  final String id;
  final String title;
  final String slug;
  final String content;
  final String excerpt;
  final String? coverImageUrl;
  final List<String> tags;
  final String authorName;
  final DateTime publishedAt;
  final String? metaTitle;
  final String? metaDescription;

  const BlogPost({
    required this.id,
    required this.title,
    required this.slug,
    required this.content,
    required this.excerpt,
    this.coverImageUrl,
    required this.tags,
    required this.authorName,
    required this.publishedAt,
    this.metaTitle,
    this.metaDescription,
  });

  factory BlogPost.fromJson(Map<String, dynamic> json) {
    return BlogPost(
      id: json['id'] as String,
      title: json['title'] as String,
      slug: json['slug'] as String,
      content: json['content'] as String,
      excerpt: json['excerpt'] as String,
      coverImageUrl: json['cover_image_url'] as String?,
      tags: (json['tags'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
      authorName: json['author_name'] as String,
      publishedAt: DateTime.parse(json['published_at'] as String),
      metaTitle: json['meta_title'] as String?,
      metaDescription: json['meta_description'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'slug': slug,
      'content': content,
      'excerpt': excerpt,
      'cover_image_url': coverImageUrl,
      'tags': tags,
      'author_name': authorName,
      'published_at': publishedAt.toIso8601String(),
      'meta_title': metaTitle,
      'meta_description': metaDescription,
    };
  }

  /// Estimates the reading time based on word count (~200 words per minute).
  String get readTime {
    final wordCount = content.split(RegExp(r'\s+')).length;
    final minutes = (wordCount / 200).ceil();
    return '$minutes min read';
  }

  /// Returns a formatted publication date (e.g., "Mar 13, 2026").
  String get formattedDate {
    const months = [
      'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
    ];
    return '${months[publishedAt.month - 1]} ${publishedAt.day}, ${publishedAt.year}';
  }
}
