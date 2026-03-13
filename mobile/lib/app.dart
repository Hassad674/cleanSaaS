import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'config/router.dart';
import 'config/theme.dart';

/// Root application widget.
///
/// Sets up [MaterialApp.router] with:
/// - GoRouter for declarative navigation with auth guards
/// - Light and dark themes matching the web design system
/// - System brightness detection for automatic theme switching
class App extends ConsumerWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = createRouter(ref);

    return MaterialApp.router(
      title: 'CleanSaaS',
      debugShowCheckedModeBanner: false,

      // Themes
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.system,

      // Navigation
      routerConfig: router,
    );
  }
}
