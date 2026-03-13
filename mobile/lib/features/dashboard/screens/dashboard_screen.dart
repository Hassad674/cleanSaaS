import 'package:flutter/material.dart';
import 'package:cleansaas_mobile/features/dashboard/widgets/stats_card.dart';
import 'package:cleansaas_mobile/features/dashboard/widgets/quick_actions.dart';
import 'package:cleansaas_mobile/features/dashboard/widgets/recent_activity.dart';

/// Main dashboard screen displaying stats, quick actions, and recent activity.
///
/// This is the home tab of the app. Uses pull-to-refresh to reload all data.
class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    if (_isLoading) return;
    setState(() => _isLoading = true);

    try {
      // TODO: Load dashboard stats from API
      await Future.delayed(const Duration(milliseconds: 500));
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Dashboard',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.person_outline),
            onPressed: () {
              // TODO: Navigate to profile
            },
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _loadData,
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  _buildGreeting(theme),
                  const SizedBox(height: 24),
                  _buildStatsSection(theme),
                  const SizedBox(height: 24),
                  const QuickActions(),
                  const SizedBox(height: 24),
                  const RecentActivity(),
                ],
              ),
      ),
    );
  }

  Widget _buildGreeting(ThemeData theme) {
    final hour = DateTime.now().hour;
    String greeting;
    if (hour < 12) {
      greeting = 'Good morning';
    } else if (hour < 17) {
      greeting = 'Good afternoon';
    } else {
      greeting = 'Good evening';
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          greeting,
          style: theme.textTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          'Here\'s what\'s happening today',
          style: theme.textTheme.bodyMedium?.copyWith(
            color: theme.colorScheme.onSurface.withOpacity(0.6),
          ),
        ),
      ],
    );
  }

  Widget _buildStatsSection(ThemeData theme) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final crossAxisCount = constraints.maxWidth > 600 ? 4 : 2;

        return GridView.count(
          crossAxisCount: crossAxisCount,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          children: [
            StatsCard(
              icon: Icons.chat_bubble_outline,
              label: 'Conversations',
              value: '12',
              color: theme.colorScheme.primary,
            ),
            StatsCard(
              icon: Icons.folder_outlined,
              label: 'Files',
              value: '34',
              color: theme.colorScheme.secondary,
            ),
            StatsCard(
              icon: Icons.notifications_outlined,
              label: 'Unread',
              value: '5',
              color: theme.colorScheme.tertiary,
            ),
            StatsCard(
              icon: Icons.storage_outlined,
              label: 'Storage',
              value: '2.4 GB',
              color: theme.colorScheme.error,
            ),
          ],
        );
      },
    );
  }
}
