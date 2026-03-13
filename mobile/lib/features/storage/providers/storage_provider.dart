import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/storage/models/file_item.dart';
import 'package:cleansaas_mobile/features/storage/repositories/storage_repository.dart';

/// State management for the file storage feature.
///
/// Manages the files list, upload progress, pagination, and error states.
class StorageProvider extends ChangeNotifier {
  final StorageRepository _repository;

  StorageProvider({required StorageRepository repository})
      : _repository = repository;

  // -- Files list state --

  List<FileItem> _files = [];
  List<FileItem> get files => _files;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _error;
  String? get error => _error;

  int _currentPage = 1;
  bool _hasMore = true;
  bool get hasMore => _hasMore;

  // -- Upload state --

  bool _isUploading = false;
  bool get isUploading => _isUploading;

  double _uploadProgress = 0.0;
  double get uploadProgress => _uploadProgress;

  String? _uploadError;
  String? get uploadError => _uploadError;

  /// Loads the first page of files, replacing any existing data.
  Future<void> loadFiles() async {
    _isLoading = true;
    _error = null;
    _currentPage = 1;
    notifyListeners();

    try {
      _files = await _repository.getFiles(page: 1);
      _hasMore = _files.length >= 20;
    } catch (e) {
      _error = 'Failed to load files. Please try again.';
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  /// Loads the next page of files and appends them to the list.
  Future<void> loadMoreFiles() async {
    if (_isLoading || !_hasMore) return;

    _isLoading = true;
    notifyListeners();

    try {
      _currentPage++;
      final more = await _repository.getFiles(page: _currentPage);
      _files = [..._files, ...more];
      _hasMore = more.length >= 20;
    } catch (e) {
      _currentPage--;
      _error = 'Failed to load more files.';
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  /// Uploads a file and adds it to the top of the list on success.
  ///
  /// Tracks upload progress via [uploadProgress] (0.0 to 1.0).
  Future<void> uploadFile({
    required String filePath,
    required String fileName,
  }) async {
    _isUploading = true;
    _uploadProgress = 0.0;
    _uploadError = null;
    notifyListeners();

    try {
      final file = await _repository.uploadFile(
        filePath: filePath,
        fileName: fileName,
        onProgress: (progress) {
          _uploadProgress = progress;
          notifyListeners();
        },
      );
      _files = [file, ..._files];
    } catch (e) {
      _uploadError = 'Failed to upload file. Please try again.';
    } finally {
      _isUploading = false;
      _uploadProgress = 0.0;
      notifyListeners();
    }
  }

  /// Deletes a file by [fileId] and removes it from the local list.
  Future<void> deleteFile(String fileId) async {
    try {
      await _repository.deleteFile(fileId);
      _files = _files.where((f) => f.id != fileId).toList();
      notifyListeners();
    } catch (e) {
      _error = 'Failed to delete file.';
      notifyListeners();
    }
  }

  /// Clears any upload error.
  void clearUploadError() {
    _uploadError = null;
    notifyListeners();
  }
}

/// Riverpod provider for [StorageProvider].
///
/// Requires a [StorageRepository] to be provided in the widget tree.
final storageProvider = ChangeNotifierProvider<StorageProvider>((ref) {
  throw UnimplementedError(
    'storageProvider must be overridden with a valid StorageRepository.',
  );
});
