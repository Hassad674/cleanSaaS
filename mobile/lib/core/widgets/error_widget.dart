import 'package:flutter/material.dart';

import '../../config/theme.dart';
import 'app_button.dart';

/// Displays an error message with an optional retry button.
///
/// Matches the web design system's destructive color for the icon.
class AppErrorWidget extends StatelessWidget {
  const AppErrorWidget({
    super.key,
    required this.message,
    this.onRetry,
    this.icon = Icons.error_outline,
  });

  final String message;
  final VoidCallback? onRetry;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final appColors = theme.extension<AppColors>();

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              icon,
              size: 48,
              color: theme.colorScheme.error,
            ),
            const SizedBox(height: 16),
            Text(
              message,
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyLarge?.copyWith(
                color: appColors?.mutedForeground ??
                    theme.colorScheme.onSurface.withOpacity(0.7),
              ),
            ),
            if (onRetry != null) ...[
              const SizedBox(height: 24),
              AppButton(
                label: 'Try again',
                onPressed: onRetry,
                isFullWidth: false,
                variant: AppButtonVariant.outlined,
                icon: Icons.refresh,
              ),
            ],
          ],
        ),
      ),
    );
  }
}
