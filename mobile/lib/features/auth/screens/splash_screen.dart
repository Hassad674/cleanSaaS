import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../config/router.dart';
import '../../../core/widgets/loading_widget.dart';
import '../providers/auth_provider.dart';

/// Initial screen displayed while the app checks for a stored session.
///
/// Shows the app logo and a loading indicator. Once [AuthNotifier] resolves
/// the stored token, [GoRouter]'s redirect logic navigates to either
/// /dashboard or /login automatically.
class SplashScreen extends ConsumerWidget {
  const SplashScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final theme = Theme.of(context);

    // The router redirect handles navigation once status is no longer loading.
    // This listener handles the edge case where the redirect fires before
    // the widget rebuilds.
    ref.listen<AuthState>(authProvider, (previous, next) {
      if (next.status == AuthStatus.authenticated) {
        context.go(RoutePaths.dashboard);
      } else if (next.status == AuthStatus.unauthenticated) {
        context.go(RoutePaths.login);
      }
    });

    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // App icon
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: theme.colorScheme.primary.withOpacity(0.1),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Icon(
                Icons.rocket_launch_rounded,
                size: 40,
                color: theme.colorScheme.primary,
              ),
            ),
            const SizedBox(height: 24),

            // App name
            Text(
              'CleanSaaS',
              style: theme.textTheme.headlineMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Loading...',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.6),
              ),
            ),
            const SizedBox(height: 40),

            // Loading spinner
            if (authState.status == AuthStatus.loading)
              const LoadingWidget(size: 28),
          ],
        ),
      ),
    );
  }
}
