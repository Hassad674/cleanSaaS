import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/shell/widgets/app_bottom_nav.dart';
import 'package:cleansaas_mobile/features/dashboard/screens/dashboard_screen.dart';
import 'package:cleansaas_mobile/features/settings/screens/settings_screen.dart';

/// Main shell wrapper providing the app scaffold with bottom navigation.
///
/// Uses [IndexedStack] to preserve the state of each tab when switching.
/// Hosts 5 tabs: Dashboard, AI, Files, Notifications, and Settings.
///
/// AI, Files, and Notifications screens are placeholders that should be
/// replaced with the actual feature screens once those features are created
/// by their respective agents.
class MainShell extends ConsumerStatefulWidget {
  const MainShell({super.key});

  @override
  ConsumerState<MainShell> createState() => _MainShellState();
}

class _MainShellState extends ConsumerState<MainShell> {
  int _currentIndex = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _currentIndex,
        children: [
          const DashboardScreen(),
          _buildPlaceholderTab(context, 'AI', Icons.auto_awesome_outlined),
          _buildPlaceholderTab(context, 'Files', Icons.folder_outlined),
          _buildPlaceholderTab(
              context, 'Notifications', Icons.notifications_outlined),
          const SettingsScreen(),
        ],
      ),
      bottomNavigationBar: AppBottomNav(
        currentIndex: _currentIndex,
        onTap: (index) {
          setState(() => _currentIndex = index);
        },
        // TODO: Connect to actual unread notification count from provider.
        unreadNotificationCount: 0,
      ),
    );
  }

  /// Builds a placeholder screen for tabs whose features have not yet been
  /// created by their respective agents.
  Widget _buildPlaceholderTab(
    BuildContext context,
    String label,
    IconData icon,
  ) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          label,
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 64,
              color: theme.colorScheme.onSurface.withOpacity(0.2),
            ),
            const SizedBox(height: 16),
            Text(
              '$label coming soon',
              style: theme.textTheme.titleMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.5),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'This feature is being built by another agent.',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.4),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
