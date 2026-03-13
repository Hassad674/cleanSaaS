import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:file_picker/file_picker.dart';
import 'package:cleansaas_mobile/features/storage/providers/storage_provider.dart';

/// A floating action button that triggers the file picker and upload flow.
///
/// Shows a bottom sheet with options to pick from gallery or files,
/// then uploads the selected file with progress tracking.
class UploadButton extends ConsumerWidget {
  const UploadButton({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final storage = ref.watch(storageProvider);

    return FloatingActionButton(
      onPressed: storage.isUploading ? null : () => _showUploadOptions(context, ref),
      child: storage.isUploading
          ? const SizedBox(
              width: 24,
              height: 24,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: Colors.white,
              ),
            )
          : const Icon(Icons.add),
    );
  }

  void _showUploadOptions(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: theme.colorScheme.onSurface.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'Upload file',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 16),
              ListTile(
                leading: Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(
                    Icons.photo_library_outlined,
                    color: theme.colorScheme.primary,
                  ),
                ),
                title: const Text('Photo or Video'),
                subtitle: const Text('Choose from your gallery'),
                onTap: () {
                  Navigator.of(ctx).pop();
                  _pickAndUpload(context, ref, FileType.media);
                },
              ),
              ListTile(
                leading: Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.secondary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(
                    Icons.insert_drive_file_outlined,
                    color: theme.colorScheme.secondary,
                  ),
                ),
                title: const Text('Document'),
                subtitle: const Text('Choose any file'),
                onTap: () {
                  Navigator.of(ctx).pop();
                  _pickAndUpload(context, ref, FileType.any);
                },
              ),
              const SizedBox(height: 8),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _pickAndUpload(
    BuildContext context,
    WidgetRef ref,
    FileType type,
  ) async {
    try {
      final result = await FilePicker.platform.pickFiles(type: type);

      if (result == null || result.files.isEmpty) return;

      final pickedFile = result.files.first;
      if (pickedFile.path == null) return;

      await ref.read(storageProvider).uploadFile(
            filePath: pickedFile.path!,
            fileName: pickedFile.name,
          );
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: const Text('Failed to pick file'),
            behavior: SnackBarBehavior.floating,
            backgroundColor: Theme.of(context).colorScheme.error,
          ),
        );
      }
    }
  }
}
