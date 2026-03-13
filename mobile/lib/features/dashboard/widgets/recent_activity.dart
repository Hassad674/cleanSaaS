import 'package:flutter/material.dart';

/// Displays a list of recent user activities on the dashboard.
///
/// Shows the latest actions across all features (chats, uploads, notifications)
/// as a timeline-style list.
class RecentActivity extends StatelessWidget {
  const RecentActivity({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Recent Activity',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            TextButton(
              onPressed: () {
                // TODO: Navigate to full activity log
              },
              child: const Text('See all'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        _buildActivityList(theme),
      ],
    );
  }

  Widget _buildActivityList(ThemeData theme) {
    // TODO: Replace with actual data from provider
    final activities = <_ActivityData>[
      _ActivityData(
        icon: Icons.chat_bubble_outline,
        title: 'New conversation started',
        subtitle: 'AI Chat',
        timeAgo: '5 min ago',
      ),
      _ActivityData(
        icon: Icons.upload_file_outlined,
        title: 'Document uploaded',
        subtitle: 'report.pdf - 2.4 MB',
        timeAgo: '1 hour ago',
      ),
      _ActivityData(
        icon: Icons.notifications_outlined,
        title: 'Subscription renewed',
        subtitle: 'Pro plan - monthly',
        timeAgo: '3 hours ago',
      ),
      _ActivityData(
        icon: Icons.security_outlined,
        title: 'Password changed',
        subtitle: 'Account security',
        timeAgo: 'Yesterday',
      ),
    ];

    if (activities.isEmpty) {
      return _buildEmptyState(theme);
    }

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: BorderSide(
          color: theme.colorScheme.outline.withOpacity(0.1),
        ),
      ),
      child: ListView.separated(
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        itemCount: activities.length,
        separatorBuilder: (context, index) => Divider(
          height: 1,
          indent: 56,
          color: theme.colorScheme.outline.withOpacity(0.1),
        ),
        itemBuilder: (context, index) {
          final activity = activities[index];
          return _ActivityTile(activity: activity);
        },
      ),
    );
  }

  Widget _buildEmptyState(ThemeData theme) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 32),
        child: Column(
          children: [
            Icon(
              Icons.history_outlined,
              size: 48,
              color: theme.colorScheme.onSurface.withOpacity(0.3),
            ),
            const SizedBox(height: 12),
            Text(
              'No recent activity',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.5),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActivityData {
  final IconData icon;
  final String title;
  final String subtitle;
  final String timeAgo;

  const _ActivityData({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.timeAgo,
  });
}

class _ActivityTile extends StatelessWidget {
  final _ActivityData activity;

  const _ActivityTile({required this.activity});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: theme.colorScheme.primary.withOpacity(0.1),
          borderRadius: BorderRadius.circular(10),
        ),
        child: Icon(
          activity.icon,
          color: theme.colorScheme.primary,
          size: 20,
        ),
      ),
      title: Text(
        activity.title,
        style: theme.textTheme.bodyMedium?.copyWith(
          fontWeight: FontWeight.w500,
        ),
      ),
      subtitle: Text(
        activity.subtitle,
        style: theme.textTheme.bodySmall?.copyWith(
          color: theme.colorScheme.onSurface.withOpacity(0.5),
        ),
      ),
      trailing: Text(
        activity.timeAgo,
        style: theme.textTheme.labelSmall?.copyWith(
          color: theme.colorScheme.onSurface.withOpacity(0.4),
        ),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
    );
  }
}
