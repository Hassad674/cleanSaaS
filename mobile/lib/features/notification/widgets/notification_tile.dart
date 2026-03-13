import 'package:flutter/material.dart';
import 'package:cleansaas_mobile/features/notification/models/notification.dart';

/// A list tile for a single notification.
///
/// Shows an unread indicator dot, icon based on notification type,
/// title, message preview, and relative timestamp.
/// Tapping marks it as read (handled by the parent).
class NotificationTile extends StatelessWidget {
  final AppNotification notification;
  final VoidCallback onTap;

  const NotificationTile({
    super.key,
    required this.notification,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isUnread = !notification.read;

    return Material(
      color: isUnread
          ? theme.colorScheme.primary.withOpacity(0.04)
          : Colors.transparent,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Icon with unread dot
              Stack(
                children: [
                  Container(
                    width: 44,
                    height: 44,
                    decoration: BoxDecoration(
                      color: _getTypeColor(theme).withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      _getTypeIcon(),
                      color: _getTypeColor(theme),
                      size: 20,
                    ),
                  ),
                  if (isUnread)
                    Positioned(
                      top: 0,
                      right: 0,
                      child: Container(
                        width: 10,
                        height: 10,
                        decoration: BoxDecoration(
                          color: theme.colorScheme.primary,
                          shape: BoxShape.circle,
                          border: Border.all(
                            color: theme.colorScheme.surface,
                            width: 2,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
              const SizedBox(width: 12),

              // Content
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            notification.title,
                            style: theme.textTheme.bodyMedium?.copyWith(
                              fontWeight:
                                  isUnread ? FontWeight.w600 : FontWeight.w400,
                            ),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                        const SizedBox(width: 8),
                        Text(
                          _formatTimeAgo(notification.createdAt),
                          style: theme.textTheme.labelSmall?.copyWith(
                            color:
                                theme.colorScheme.onSurface.withOpacity(0.4),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(
                      notification.message,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurface.withOpacity(
                          isUnread ? 0.7 : 0.5,
                        ),
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  IconData _getTypeIcon() {
    switch (notification.type) {
      case 'billing':
        return Icons.payment_outlined;
      case 'security':
        return Icons.security_outlined;
      case 'ai':
        return Icons.smart_toy_outlined;
      case 'storage':
        return Icons.cloud_outlined;
      case 'warning':
        return Icons.warning_amber_outlined;
      case 'error':
        return Icons.error_outline;
      case 'success':
        return Icons.check_circle_outline;
      case 'info':
      default:
        return Icons.info_outline;
    }
  }

  Color _getTypeColor(ThemeData theme) {
    switch (notification.type) {
      case 'billing':
        return theme.colorScheme.secondary;
      case 'security':
        return theme.colorScheme.error;
      case 'ai':
        return theme.colorScheme.primary;
      case 'storage':
        return theme.colorScheme.tertiary;
      case 'warning':
        return theme.colorScheme.tertiary;
      case 'error':
        return theme.colorScheme.error;
      case 'success':
        return theme.colorScheme.primary;
      case 'info':
      default:
        return theme.colorScheme.primary;
    }
  }

  String _formatTimeAgo(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inSeconds < 60) return 'Just now';
    if (difference.inMinutes < 60) return '${difference.inMinutes}m';
    if (difference.inHours < 24) return '${difference.inHours}h';
    if (difference.inDays < 7) return '${difference.inDays}d';
    if (difference.inDays < 30) return '${(difference.inDays / 7).floor()}w';

    return '${dateTime.day}/${dateTime.month}';
  }
}
