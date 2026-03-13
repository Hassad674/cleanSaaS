import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/notification/providers/notification_provider.dart';
import 'package:cleansaas_mobile/features/notification/widgets/notification_tile.dart';

/// Screen displaying the user's notifications with unread filter and mark-all-read.
///
/// Supports pull-to-refresh, infinite scroll pagination, an unread-only filter
/// toggle, and a "Mark all as read" action in the app bar.
class NotificationsScreen extends ConsumerStatefulWidget {
  const NotificationsScreen({super.key});

  @override
  ConsumerState<NotificationsScreen> createState() =>
      _NotificationsScreenState();
}

class _NotificationsScreenState extends ConsumerState<NotificationsScreen> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(notificationProvider).loadNotifications();
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
      ref.read(notificationProvider).loadMoreNotifications();
    }
  }

  Future<void> _onRefresh() async {
    await ref.read(notificationProvider).loadNotifications();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final notifications = ref.watch(notificationProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Notifications',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        actions: [
          if (notifications.unreadCount > 0)
            TextButton(
              onPressed: () => _confirmMarkAllRead(context, notifications),
              child: const Text('Mark all read'),
            ),
        ],
      ),
      body: Column(
        children: [
          _buildFilterBar(theme, notifications),
          Expanded(child: _buildBody(theme, notifications)),
        ],
      ),
    );
  }

  Widget _buildFilterBar(ThemeData theme, NotificationProvider notifications) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        border: Border(
          bottom: BorderSide(
            color: theme.colorScheme.outline.withOpacity(0.1),
          ),
        ),
      ),
      child: Row(
        children: [
          FilterChip(
            label: const Text('All'),
            selected: !notifications.showUnreadOnly,
            onSelected: (selected) {
              if (selected && notifications.showUnreadOnly) {
                notifications.toggleUnreadFilter();
              }
            },
          ),
          const SizedBox(width: 8),
          FilterChip(
            label: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text('Unread'),
                if (notifications.unreadCount > 0) ...[
                  const SizedBox(width: 6),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 6,
                      vertical: 1,
                    ),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primary,
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: Text(
                      '${notifications.unreadCount}',
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: theme.colorScheme.onPrimary,
                        fontSize: 10,
                      ),
                    ),
                  ),
                ],
              ],
            ),
            selected: notifications.showUnreadOnly,
            onSelected: (selected) {
              if (selected && !notifications.showUnreadOnly) {
                notifications.toggleUnreadFilter();
              }
            },
          ),
        ],
      ),
    );
  }

  Widget _buildBody(ThemeData theme, NotificationProvider notifications) {
    if (notifications.isLoading && notifications.notifications.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (notifications.error != null && notifications.notifications.isEmpty) {
      return _buildErrorState(theme, notifications);
    }

    if (notifications.notifications.isEmpty) {
      return _buildEmptyState(theme, notifications);
    }

    return RefreshIndicator(
      onRefresh: _onRefresh,
      child: ListView.separated(
        controller: _scrollController,
        padding: const EdgeInsets.symmetric(vertical: 4),
        itemCount: notifications.notifications.length +
            (notifications.hasMore ? 1 : 0),
        separatorBuilder: (context, index) => Divider(
          height: 1,
          indent: 72,
          color: theme.colorScheme.outline.withOpacity(0.1),
        ),
        itemBuilder: (context, index) {
          if (index >= notifications.notifications.length) {
            return const Padding(
              padding: EdgeInsets.all(16),
              child: Center(child: CircularProgressIndicator()),
            );
          }

          final notification = notifications.notifications[index];
          return NotificationTile(
            notification: notification,
            onTap: () {
              if (!notification.read) {
                notifications.markAsRead(notification.id);
              }
            },
          );
        },
      ),
    );
  }

  void _confirmMarkAllRead(
    BuildContext context,
    NotificationProvider notifications,
  ) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Mark all as read?'),
        content: Text(
          'This will mark ${notifications.unreadCount} notification${notifications.unreadCount == 1 ? '' : 's'} as read.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              Navigator.of(ctx).pop();
              notifications.markAllAsRead();
            },
            child: const Text('Mark all read'),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(
    ThemeData theme,
    NotificationProvider notifications,
  ) {
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
              notifications.error!,
              style: theme.textTheme.bodyLarge,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: () => notifications.loadNotifications(),
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(
    ThemeData theme,
    NotificationProvider notifications,
  ) {
    final message = notifications.showUnreadOnly
        ? 'No unread notifications'
        : 'No notifications yet';
    final subtitle = notifications.showUnreadOnly
        ? 'All caught up!'
        : 'We\'ll notify you when something happens';

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.notifications_none_outlined,
              size: 64,
              color: theme.colorScheme.onSurface.withOpacity(0.3),
            ),
            const SizedBox(height: 16),
            Text(
              message,
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              subtitle,
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
