import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/settings/providers/settings_provider.dart';
import 'package:cleansaas_mobile/features/settings/widgets/profile_header.dart';
import 'package:cleansaas_mobile/features/settings/widgets/settings_tile.dart';
import 'package:cleansaas_mobile/features/settings/widgets/danger_zone.dart';
import 'package:cleansaas_mobile/features/settings/screens/edit_profile_screen.dart';
import 'package:cleansaas_mobile/features/settings/screens/change_password_screen.dart';

/// Main settings screen with profile header and organized setting sections.
///
/// Sections:
/// - Profile header (avatar + name + email)
/// - Account (edit profile, change password)
/// - Preferences (dark mode toggle, notifications toggle)
/// - Billing (view plan, manage subscription)
/// - Danger Zone (delete account)
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final profileAsync = ref.watch(userProfileProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Settings',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: profileAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => _buildErrorState(context, ref, error),
        data: (profile) => RefreshIndicator(
          onRefresh: () async {
            ref.invalidate(userProfileProvider);
          },
          child: ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Profile header
              ProfileHeader(
                profile: profile,
                onTap: () => _navigateToEditProfile(context),
              ),
              const SizedBox(height: 24),

              // Account section
              _buildSectionHeader(theme, 'Account'),
              const SizedBox(height: 8),
              _buildSectionCard(
                theme: theme,
                children: [
                  SettingsTile(
                    icon: Icons.person_outline,
                    label: 'Edit Profile',
                    subtitle: 'Name, email, avatar',
                    onTap: () => _navigateToEditProfile(context),
                  ),
                  Divider(
                    height: 1,
                    indent: 70,
                    color: theme.colorScheme.outlineVariant,
                  ),
                  SettingsTile(
                    icon: Icons.lock_outline,
                    label: 'Change Password',
                    subtitle: 'Update your password',
                    onTap: () => _navigateToChangePassword(context),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // Preferences section
              _buildSectionHeader(theme, 'Preferences'),
              const SizedBox(height: 8),
              _buildSectionCard(
                theme: theme,
                children: [
                  SettingsTile(
                    icon: Icons.dark_mode_outlined,
                    label: 'Dark Mode',
                    subtitle: _darkModeSubtitle(ref),
                    trailing: Switch.adaptive(
                      value: ref.watch(darkModeProvider) ?? false,
                      onChanged: (value) {
                        ref.read(darkModeProvider.notifier).state = value;
                      },
                    ),
                  ),
                  Divider(
                    height: 1,
                    indent: 70,
                    color: theme.colorScheme.outlineVariant,
                  ),
                  SettingsTile(
                    icon: Icons.notifications_outlined,
                    label: 'Notifications',
                    subtitle: ref.watch(notificationsEnabledProvider)
                        ? 'Enabled'
                        : 'Disabled',
                    trailing: Switch.adaptive(
                      value: ref.watch(notificationsEnabledProvider),
                      onChanged: (value) {
                        ref.read(notificationsEnabledProvider.notifier).state =
                            value;
                      },
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // Billing section
              _buildSectionHeader(theme, 'Billing'),
              const SizedBox(height: 8),
              _buildSectionCard(
                theme: theme,
                children: [
                  SettingsTile(
                    icon: Icons.credit_card_outlined,
                    label: 'Current Plan',
                    subtitle: 'View your subscription',
                    onTap: () {
                      // Navigate to billing screen via shell navigation
                      // This is handled at the app level since billing is
                      // a separate feature module.
                    },
                  ),
                  Divider(
                    height: 1,
                    indent: 70,
                    color: theme.colorScheme.outlineVariant,
                  ),
                  SettingsTile(
                    icon: Icons.receipt_long_outlined,
                    label: 'Manage Subscription',
                    subtitle: 'Upgrade, downgrade, or cancel',
                    onTap: () {
                      // Navigate to plans screen
                    },
                  ),
                ],
              ),
              const SizedBox(height: 32),

              // Danger zone
              DangerZone(
                onDeleteAccount: (password) async {
                  await ref
                      .read(userProfileProvider.notifier)
                      .deleteAccount(password: password);
                  // After deletion, the auth layer should handle logout/redirect.
                },
              ),
              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
    );
  }

  String _darkModeSubtitle(WidgetRef ref) {
    final darkMode = ref.watch(darkModeProvider);
    if (darkMode == null) return 'System default';
    return darkMode ? 'On' : 'Off';
  }

  Widget _buildErrorState(BuildContext context, WidgetRef ref, Object error) {
    final theme = Theme.of(context);
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.error_outline,
              size: 48,
              color: theme.colorScheme.error,
            ),
            const SizedBox(height: 16),
            Text(
              'Failed to load profile',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              error.toString(),
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.6),
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: () => ref.invalidate(userProfileProvider),
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(ThemeData theme, String title) {
    return Padding(
      padding: const EdgeInsets.only(left: 4),
      child: Text(
        title,
        style: theme.textTheme.titleSmall?.copyWith(
          fontWeight: FontWeight.bold,
          color: theme.colorScheme.onSurface.withOpacity(0.5),
          letterSpacing: 0.5,
        ),
      ),
    );
  }

  Widget _buildSectionCard({
    required ThemeData theme,
    required List<Widget> children,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
      ),
      clipBehavior: Clip.antiAlias,
      child: Column(children: children),
    );
  }

  void _navigateToEditProfile(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => const EditProfileScreen(),
      ),
    );
  }

  void _navigateToChangePassword(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => const ChangePasswordScreen(),
      ),
    );
  }
}
