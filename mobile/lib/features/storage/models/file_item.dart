/// Represents a file stored in the system.
///
/// Maps to the backend `files` table. Contains metadata about uploaded files
/// including name, URL, size, and content type.
class FileItem {
  final String id;
  final String name;
  final String url;
  final int size;
  final String contentType;
  final DateTime createdAt;

  const FileItem({
    required this.id,
    required this.name,
    required this.url,
    required this.size,
    required this.contentType,
    required this.createdAt,
  });

  factory FileItem.fromJson(Map<String, dynamic> json) {
    return FileItem(
      id: json['id'] as String,
      name: json['name'] as String,
      url: json['url'] as String,
      size: json['size'] as int,
      contentType: json['content_type'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'url': url,
      'size': size,
      'content_type': contentType,
      'created_at': createdAt.toIso8601String(),
    };
  }

  /// Returns a human-readable file size string.
  String get formattedSize {
    if (size < 1024) return '$size B';
    if (size < 1024 * 1024) return '${(size / 1024).toStringAsFixed(1)} KB';
    if (size < 1024 * 1024 * 1024) {
      return '${(size / (1024 * 1024)).toStringAsFixed(1)} MB';
    }
    return '${(size / (1024 * 1024 * 1024)).toStringAsFixed(1)} GB';
  }

  /// Returns the file extension (e.g. "pdf", "jpg").
  String get extension {
    final parts = name.split('.');
    return parts.length > 1 ? parts.last.toLowerCase() : '';
  }

  /// Whether this file is an image based on its content type.
  bool get isImage => contentType.startsWith('image/');

  /// Whether this file is a video based on its content type.
  bool get isVideo => contentType.startsWith('video/');

  /// Whether this file is a PDF.
  bool get isPdf => contentType == 'application/pdf';

  FileItem copyWith({
    String? id,
    String? name,
    String? url,
    int? size,
    String? contentType,
    DateTime? createdAt,
  }) {
    return FileItem(
      id: id ?? this.id,
      name: name ?? this.name,
      url: url ?? this.url,
      size: size ?? this.size,
      contentType: contentType ?? this.contentType,
      createdAt: createdAt ?? this.createdAt,
    );
  }
}
