import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../features/ai/screens/chat_screen.dart';
import '../features/ai/screens/conversations_screen.dart';
import '../features/auth/providers/auth_provider.dart';
import '../features/auth/screens/forgot_password_screen.dart';
import '../features/auth/screens/login_screen.dart';
import '../features/auth/screens/register_screen.dart';
import '../features/auth/screens/splash_screen.dart';
import '../features/billing/screens/billing_screen.dart';
import '../features/dashboard/screens/dashboard_screen.dart';
import '../features/settings/screens/settings_screen.dart';
import '../features/storage/screens/files_screen.dart';

/// Route path constants to avoid magic strings.
class RoutePaths {
  RoutePaths._();

  static const String splash = '/splash';
  static const String login = '/login';
  static const String register = '/register';
  static const String forgotPassword = '/forgot-password';
  static const String dashboard = '/dashboard';
  static const String ai = '/ai';
  static const String aiChat = '/ai/:id';
  static const String files = '/files';
  static const String notifications = '/notifications';
  static const String settings = '/settings';
  static const String billing = '/settings/billing';
}

/// Creates the app router with authentication-based redirects.
///
/// Watches [authProvider] to determine whether the user is authenticated and
/// redirects to /login or /dashboard accordingly.
GoRouter createRouter(WidgetRef ref) {
  return GoRouter(
    initialLocation: RoutePaths.splash,
    debugLogDiagnostics: true,
    redirect: (BuildContext context, GoRouterState state) {
      final authState = ref.read(authProvider);
      final isAuthRoute = state.matchedLocation == RoutePaths.login ||
          state.matchedLocation == RoutePaths.register ||
          state.matchedLocation == RoutePaths.forgotPassword ||
          state.matchedLocation == RoutePaths.splash;

      // Still loading — let splash screen handle it
      if (authState.status == AuthStatus.loading) {
        return state.matchedLocation == RoutePaths.splash
            ? null
            : RoutePaths.splash;
      }

      // Not authenticated — force to login (unless already on an auth route)
      if (authState.status == AuthStatus.unauthenticated) {
        return isAuthRoute ? null : RoutePaths.login;
      }

      // Authenticated — redirect away from auth routes
      if (authState.status == AuthStatus.authenticated && isAuthRoute) {
        return RoutePaths.dashboard;
      }

      return null;
    },
    routes: [
      // --- Auth routes ---
      GoRoute(
        path: RoutePaths.splash,
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: RoutePaths.login,
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: RoutePaths.register,
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: RoutePaths.forgotPassword,
        builder: (context, state) => const ForgotPasswordScreen(),
      ),

      // --- Authenticated routes ---
      GoRoute(
        path: RoutePaths.dashboard,
        builder: (context, state) => const DashboardScreen(),
      ),
      GoRoute(
        path: RoutePaths.ai,
        builder: (context, state) => const ConversationsScreen(),
        routes: [
          GoRoute(
            path: ':id',
            builder: (context, state) {
              final conversationId = state.pathParameters['id']!;
              return ChatScreen(conversationId: conversationId);
            },
          ),
        ],
      ),
      GoRoute(
        path: RoutePaths.files,
        builder: (context, state) => const FilesScreen(),
      ),
      GoRoute(
        path: RoutePaths.notifications,
        builder: (context, state) =>
            const _PlaceholderScreen(title: 'Notifications'),
      ),
      GoRoute(
        path: RoutePaths.settings,
        builder: (context, state) => const SettingsScreen(),
        routes: [
          GoRoute(
            path: 'billing',
            builder: (context, state) => const BillingScreen(),
          ),
        ],
      ),
    ],
  );
}

/// Temporary placeholder for screens that have not been implemented yet.
///
/// Replace with the real feature screen when it is built.
class _PlaceholderScreen extends StatelessWidget {
  const _PlaceholderScreen({required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: Text(title)),
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.construction, size: 48, color: theme.colorScheme.primary),
            const SizedBox(height: 16),
            Text(
              title,
              style: theme.textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            Text(
              'Coming soon',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.6),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
