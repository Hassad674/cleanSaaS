import 'package:dio/dio.dart';
import 'package:cleansaas_mobile/core/api/api_client.dart';
import 'package:cleansaas_mobile/features/storage/models/file_item.dart';

/// Repository handling all file storage API communication.
///
/// Provides methods for listing, uploading, and deleting files.
/// Upload uses multipart form data through the [ApiClient].
class StorageRepository {
  final ApiClient _apiClient;

  StorageRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Fetches a paginated list of files for the current user.
  Future<List<FileItem>> getFiles({
    int page = 1,
    int limit = 20,
  }) async {
    final response = await _apiClient.get(
      '/files',
      queryParameters: {
        'page': page.toString(),
        'limit': limit.toString(),
      },
    );

    final List<dynamic> data = response.data['data'] as List<dynamic>;
    return data
        .map((json) => FileItem.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Uploads a file using multipart form data.
  ///
  /// [filePath] is the local filesystem path to the file.
  /// [fileName] is the desired name for the uploaded file.
  /// [onProgress] is an optional callback receiving upload progress (0.0 - 1.0).
  Future<FileItem> uploadFile({
    required String filePath,
    required String fileName,
    void Function(double progress)? onProgress,
  }) async {
    final formData = FormData.fromMap({
      'file': await MultipartFile.fromFile(filePath, filename: fileName),
    });

    final response = await _apiClient.post(
      '/files/upload',
      data: formData,
    );

    return FileItem.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Deletes a file by its [fileId].
  Future<void> deleteFile(String fileId) async {
    await _apiClient.delete('/files/$fileId');
  }
}
