import 'package:flutter/material.dart';

/// Bottom navigation bar for the main app shell.
///
/// Displays 5 tabs: Dashboard, AI, Files, Notifications, and Settings.
/// Supports a badge on the notifications tab showing the unread count.
class AppBottomNav extends StatelessWidget {
  final int currentIndex;
  final ValueChanged<int> onTap;
  final int unreadNotificationCount;

  const AppBottomNav({
    super.key,
    required this.currentIndex,
    required this.onTap,
    this.unreadNotificationCount = 0,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return NavigationBar(
      selectedIndex: currentIndex,
      onDestinationSelected: onTap,
      labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
      backgroundColor: theme.colorScheme.surface,
      indicatorColor: theme.colorScheme.primaryContainer,
      destinations: [
        const NavigationDestination(
          icon: Icon(Icons.dashboard_outlined),
          selectedIcon: Icon(Icons.dashboard),
          label: 'Dashboard',
        ),
        const NavigationDestination(
          icon: Icon(Icons.auto_awesome_outlined),
          selectedIcon: Icon(Icons.auto_awesome),
          label: 'AI',
        ),
        const NavigationDestination(
          icon: Icon(Icons.folder_outlined),
          selectedIcon: Icon(Icons.folder),
          label: 'Files',
        ),
        NavigationDestination(
          icon: _buildNotificationIcon(
            theme,
            isSelected: false,
          ),
          selectedIcon: _buildNotificationIcon(
            theme,
            isSelected: true,
          ),
          label: 'Notifications',
        ),
        const NavigationDestination(
          icon: Icon(Icons.settings_outlined),
          selectedIcon: Icon(Icons.settings),
          label: 'Settings',
        ),
      ],
    );
  }

  Widget _buildNotificationIcon(
    ThemeData theme, {
    required bool isSelected,
  }) {
    final icon = Icon(
      isSelected ? Icons.notifications : Icons.notifications_outlined,
    );

    if (unreadNotificationCount <= 0) return icon;

    return Badge(
      label: Text(
        unreadNotificationCount > 99 ? '99+' : '$unreadNotificationCount',
        style: const TextStyle(fontSize: 10),
      ),
      child: icon,
    );
  }
}
