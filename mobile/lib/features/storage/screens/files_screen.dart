import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/storage/providers/storage_provider.dart';
import 'package:cleansaas_mobile/features/storage/widgets/file_tile.dart';
import 'package:cleansaas_mobile/features/storage/widgets/upload_button.dart';

/// Screen displaying the user's uploaded files with grid/list toggle.
///
/// Supports pull-to-refresh, infinite scroll pagination, view mode toggle
/// (grid vs list), and a floating action button for file uploads.
class FilesScreen extends ConsumerStatefulWidget {
  const FilesScreen({super.key});

  @override
  ConsumerState<FilesScreen> createState() => _FilesScreenState();
}

class _FilesScreenState extends ConsumerState<FilesScreen> {
  final ScrollController _scrollController = ScrollController();
  bool _isGridView = true;

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(storageProvider).loadFiles();
    });
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
      ref.read(storageProvider).loadMoreFiles();
    }
  }

  Future<void> _onRefresh() async {
    await ref.read(storageProvider).loadFiles();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final storage = ref.watch(storageProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Files',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        actions: [
          IconButton(
            icon: Icon(_isGridView ? Icons.view_list : Icons.grid_view),
            onPressed: () {
              setState(() => _isGridView = !_isGridView);
            },
            tooltip: _isGridView ? 'List view' : 'Grid view',
          ),
        ],
      ),
      floatingActionButton: const UploadButton(),
      body: Column(
        children: [
          if (storage.isUploading) _buildUploadProgress(theme, storage),
          if (storage.uploadError != null)
            _buildUploadError(theme, storage),
          Expanded(child: _buildBody(theme, storage)),
        ],
      ),
    );
  }

  Widget _buildUploadProgress(ThemeData theme, StorageProvider storage) {
    return Container(
      padding: const EdgeInsets.all(16),
      color: theme.colorScheme.primaryContainer.withOpacity(0.3),
      child: Row(
        children: [
          const SizedBox(
            width: 20,
            height: 20,
            child: CircularProgressIndicator(strokeWidth: 2),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Uploading...',
                  style: theme.textTheme.bodySmall?.copyWith(
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 4),
                ClipRRect(
                  borderRadius: BorderRadius.circular(4),
                  child: LinearProgressIndicator(
                    value: storage.uploadProgress,
                    minHeight: 4,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 12),
          Text(
            '${(storage.uploadProgress * 100).toInt()}%',
            style: theme.textTheme.labelSmall,
          ),
        ],
      ),
    );
  }

  Widget _buildUploadError(ThemeData theme, StorageProvider storage) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: theme.colorScheme.errorContainer,
      child: Row(
        children: [
          Icon(
            Icons.error_outline,
            size: 16,
            color: theme.colorScheme.onErrorContainer,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              storage.uploadError!,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onErrorContainer,
              ),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.close, size: 16),
            onPressed: () => storage.clearUploadError(),
            padding: EdgeInsets.zero,
            constraints: const BoxConstraints(),
          ),
        ],
      ),
    );
  }

  Widget _buildBody(ThemeData theme, StorageProvider storage) {
    if (storage.isLoading && storage.files.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (storage.error != null && storage.files.isEmpty) {
      return _buildErrorState(theme, storage);
    }

    if (storage.files.isEmpty) {
      return _buildEmptyState(theme);
    }

    return RefreshIndicator(
      onRefresh: _onRefresh,
      child: _isGridView ? _buildGridView(storage) : _buildListView(storage),
    );
  }

  Widget _buildGridView(StorageProvider storage) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final crossAxisCount = constraints.maxWidth > 600 ? 4 : 2;

        return GridView.builder(
          controller: _scrollController,
          padding: const EdgeInsets.all(16),
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: crossAxisCount,
            crossAxisSpacing: 12,
            mainAxisSpacing: 12,
            childAspectRatio: 0.85,
          ),
          itemCount: storage.files.length + (storage.hasMore ? 1 : 0),
          itemBuilder: (context, index) {
            if (index >= storage.files.length) {
              return const Center(child: CircularProgressIndicator());
            }
            return FileTile(
              file: storage.files[index],
              isGrid: true,
              onDelete: () => _confirmDelete(context, storage, index),
            );
          },
        );
      },
    );
  }

  Widget _buildListView(StorageProvider storage) {
    return ListView.separated(
      controller: _scrollController,
      padding: const EdgeInsets.symmetric(vertical: 8),
      itemCount: storage.files.length + (storage.hasMore ? 1 : 0),
      separatorBuilder: (context, index) => Divider(
        height: 1,
        indent: 72,
        color: Theme.of(context).colorScheme.outline.withOpacity(0.1),
      ),
      itemBuilder: (context, index) {
        if (index >= storage.files.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: Center(child: CircularProgressIndicator()),
          );
        }
        return FileTile(
          file: storage.files[index],
          isGrid: false,
          onDelete: () => _confirmDelete(context, storage, index),
        );
      },
    );
  }

  void _confirmDelete(
    BuildContext context,
    StorageProvider storage,
    int index,
  ) {
    final file = storage.files[index];
    final theme = Theme.of(context);

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Delete file?'),
        content: Text('Are you sure you want to delete "${file.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              Navigator.of(ctx).pop();
              storage.deleteFile(file.id);
            },
            style: FilledButton.styleFrom(
              backgroundColor: theme.colorScheme.error,
            ),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(ThemeData theme, StorageProvider storage) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.error_outline,
              size: 48,
              color: theme.colorScheme.error,
            ),
            const SizedBox(height: 16),
            Text(
              storage.error!,
              style: theme.textTheme.bodyLarge,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: () => storage.loadFiles(),
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(ThemeData theme) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.folder_outlined,
              size: 64,
              color: theme.colorScheme.onSurface.withOpacity(0.3),
            ),
            const SizedBox(height: 16),
            Text(
              'No files yet',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Upload your first file to get started',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.5),
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}
