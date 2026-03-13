import 'package:flutter/material.dart';
import 'package:cleansaas_mobile/features/storage/models/file_item.dart';

/// Displays a single file as either a grid card or list tile.
///
/// Shows a file-type icon, name, size, and date. Supports tap for preview
/// and a delete action.
class FileTile extends StatelessWidget {
  final FileItem file;
  final bool isGrid;
  final VoidCallback onDelete;

  const FileTile({
    super.key,
    required this.file,
    required this.isGrid,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return isGrid ? _buildGridTile(context) : _buildListTile(context);
  }

  Widget _buildGridTile(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(
          color: theme.colorScheme.outline.withOpacity(0.1),
        ),
      ),
      child: InkWell(
        onTap: () {
          // TODO: Open file preview
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  _FileIcon(file: file, size: 40),
                  PopupMenuButton<String>(
                    onSelected: (value) {
                      if (value == 'delete') onDelete();
                    },
                    itemBuilder: (ctx) => [
                      const PopupMenuItem(
                        value: 'delete',
                        child: Row(
                          children: [
                            Icon(Icons.delete_outline, size: 18),
                            SizedBox(width: 8),
                            Text('Delete'),
                          ],
                        ),
                      ),
                    ],
                    icon: Icon(
                      Icons.more_vert,
                      size: 18,
                      color: theme.colorScheme.onSurface.withOpacity(0.4),
                    ),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                  ),
                ],
              ),
              const Spacer(),
              Text(
                file.name,
                style: theme.textTheme.bodySmall?.copyWith(
                  fontWeight: FontWeight.w500,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),
              Text(
                file.formattedSize,
                style: theme.textTheme.labelSmall?.copyWith(
                  color: theme.colorScheme.onSurface.withOpacity(0.5),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildListTile(BuildContext context) {
    final theme = Theme.of(context);

    return ListTile(
      onTap: () {
        // TODO: Open file preview
      },
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      leading: _FileIcon(file: file, size: 44),
      title: Text(
        file.name,
        style: theme.textTheme.bodyMedium?.copyWith(
          fontWeight: FontWeight.w500,
        ),
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      ),
      subtitle: Text(
        '${file.formattedSize} - ${_formatDate(file.createdAt)}',
        style: theme.textTheme.bodySmall?.copyWith(
          color: theme.colorScheme.onSurface.withOpacity(0.5),
        ),
      ),
      trailing: IconButton(
        icon: Icon(
          Icons.delete_outline,
          color: theme.colorScheme.error.withOpacity(0.7),
          size: 20,
        ),
        onPressed: onDelete,
      ),
    );
  }

  String _formatDate(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inDays == 0) return 'Today';
    if (difference.inDays == 1) return 'Yesterday';
    if (difference.inDays < 7) return '${difference.inDays}d ago';

    return '${dateTime.day}/${dateTime.month}/${dateTime.year}';
  }
}

/// Renders a file-type appropriate icon with a colored background.
class _FileIcon extends StatelessWidget {
  final FileItem file;
  final double size;

  const _FileIcon({required this.file, required this.size});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final iconData = _getIconForFile(file);
    final color = _getColorForFile(file, theme);

    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(size * 0.25),
      ),
      child: Icon(
        iconData,
        color: color,
        size: size * 0.5,
      ),
    );
  }

  IconData _getIconForFile(FileItem file) {
    if (file.isImage) return Icons.image_outlined;
    if (file.isVideo) return Icons.videocam_outlined;
    if (file.isPdf) return Icons.picture_as_pdf_outlined;

    switch (file.extension) {
      case 'doc':
      case 'docx':
        return Icons.description_outlined;
      case 'xls':
      case 'xlsx':
        return Icons.table_chart_outlined;
      case 'ppt':
      case 'pptx':
        return Icons.slideshow_outlined;
      case 'zip':
      case 'rar':
      case '7z':
        return Icons.folder_zip_outlined;
      case 'mp3':
      case 'wav':
      case 'aac':
        return Icons.audio_file_outlined;
      case 'txt':
      case 'md':
        return Icons.text_snippet_outlined;
      case 'json':
      case 'xml':
      case 'csv':
        return Icons.data_object_outlined;
      default:
        return Icons.insert_drive_file_outlined;
    }
  }

  Color _getColorForFile(FileItem file, ThemeData theme) {
    if (file.isImage) return theme.colorScheme.tertiary;
    if (file.isVideo) return theme.colorScheme.error;
    if (file.isPdf) return theme.colorScheme.error;

    switch (file.extension) {
      case 'doc':
      case 'docx':
        return theme.colorScheme.primary;
      case 'xls':
      case 'xlsx':
        return theme.colorScheme.secondary;
      case 'ppt':
      case 'pptx':
        return theme.colorScheme.tertiary;
      default:
        return theme.colorScheme.primary;
    }
  }
}
