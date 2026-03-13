import 'package:cleansaas_mobile/core/api/api_client.dart';
import 'package:cleansaas_mobile/features/blog/models/blog_post.dart';

/// Repository handling all blog-related API operations.
///
/// Provides methods to fetch posts (with optional tag filtering),
/// fetch a single post by slug, and list all available tags.
class BlogRepository {
  final ApiClient _apiClient;

  BlogRepository(this._apiClient);

  /// Fetches a paginated list of blog posts.
  ///
  /// Optionally filters by [tag]. Supports pagination via [page] and [limit].
  Future<List<BlogPost>> getPosts({
    String? tag,
    int page = 1,
    int limit = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'limit': limit,
    };
    if (tag != null && tag.isNotEmpty) {
      queryParams['tag'] = tag;
    }

    final response = await _apiClient.get(
      '/blog/posts',
      queryParameters: queryParams,
    );

    final posts = (response.data as List<dynamic>)
        .map((json) => BlogPost.fromJson(json as Map<String, dynamic>))
        .toList();
    return posts;
  }

  /// Fetches a single blog post by its slug.
  Future<BlogPost> getPostBySlug(String slug) async {
    final response = await _apiClient.get('/blog/posts/$slug');
    return BlogPost.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches all available blog tags.
  Future<List<String>> getTags() async {
    final response = await _apiClient.get('/blog/tags');
    final tags = (response.data as List<dynamic>)
        .map((e) => e as String)
        .toList();
    return tags;
  }
}
