import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:cleansaas_mobile/features/billing/models/subscription.dart';

/// Displays the current subscription status with plan name, status badge,
/// renewal date, and a cancel/resume button.
class SubscriptionStatus extends StatelessWidget {
  final Subscription subscription;
  final VoidCallback? onCancel;
  final VoidCallback? onResume;
  final VoidCallback? onChangePlan;

  const SubscriptionStatus({
    super.key,
    required this.subscription,
    this.onCancel,
    this.onResume,
    this.onChangePlan,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final dateFormat = DateFormat('MMM d, yyyy');

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header row: plan name + status badge
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                subscription.planName,
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              _buildStatusBadge(theme),
            ],
          ),
          const SizedBox(height: 16),

          // Renewal / expiry date
          _buildInfoRow(
            theme,
            icon: Icons.calendar_today_outlined,
            label: subscription.isCanceledButActive
                ? 'Expires'
                : 'Renews',
            value: dateFormat.format(subscription.currentPeriodEnd),
          ),
          const SizedBox(height: 8),

          // Member since
          _buildInfoRow(
            theme,
            icon: Icons.access_time,
            label: 'Member since',
            value: dateFormat.format(subscription.createdAt),
          ),
          const SizedBox(height: 20),

          // Action buttons
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: onChangePlan,
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: const Text('Change Plan'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: subscription.isCanceledButActive
                    ? FilledButton(
                        onPressed: onResume,
                        style: FilledButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: const Text('Resume'),
                      )
                    : OutlinedButton(
                        onPressed: onCancel,
                        style: OutlinedButton.styleFrom(
                          foregroundColor: theme.colorScheme.error,
                          side: BorderSide(color: theme.colorScheme.error),
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: const Text('Cancel'),
                      ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStatusBadge(ThemeData theme) {
    Color backgroundColor;
    Color foregroundColor;

    switch (subscription.status) {
      case 'active':
        backgroundColor = Colors.green.withOpacity(0.1);
        foregroundColor = Colors.green;
        break;
      case 'trialing':
        backgroundColor = theme.colorScheme.primary.withOpacity(0.1);
        foregroundColor = theme.colorScheme.primary;
        break;
      case 'canceled':
        backgroundColor = theme.colorScheme.error.withOpacity(0.1);
        foregroundColor = theme.colorScheme.error;
        break;
      case 'past_due':
        backgroundColor = Colors.orange.withOpacity(0.1);
        foregroundColor = Colors.orange;
        break;
      default:
        backgroundColor = theme.colorScheme.surfaceContainerHighest;
        foregroundColor = theme.colorScheme.onSurface;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        subscription.statusLabel,
        style: theme.textTheme.labelSmall?.copyWith(
          color: foregroundColor,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  Widget _buildInfoRow(
    ThemeData theme, {
    required IconData icon,
    required String label,
    required String value,
  }) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: theme.colorScheme.onSurface.withOpacity(0.5),
        ),
        const SizedBox(width: 8),
        Text(
          '$label: ',
          style: theme.textTheme.bodyMedium?.copyWith(
            color: theme.colorScheme.onSurface.withOpacity(0.5),
          ),
        ),
        Text(
          value,
          style: theme.textTheme.bodyMedium?.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}
