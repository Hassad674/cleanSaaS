import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/billing/providers/billing_provider.dart';
import 'package:cleansaas_mobile/features/billing/models/subscription.dart';
import 'package:cleansaas_mobile/features/billing/widgets/subscription_status.dart';
import 'package:cleansaas_mobile/features/billing/widgets/invoice_tile.dart';
import 'package:cleansaas_mobile/features/billing/screens/plans_screen.dart';

/// Main billing screen displaying the current subscription status,
/// plan details, invoice history, and management actions.
class BillingScreen extends ConsumerWidget {
  const BillingScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final subscriptionAsync = ref.watch(subscriptionProvider);
    final invoicesAsync = ref.watch(invoicesProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Billing',
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(subscriptionProvider);
          ref.invalidate(invoicesProvider);
        },
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Subscription section
            subscriptionAsync.when(
              loading: () => const _SubscriptionSkeleton(),
              error: (error, _) => _buildErrorCard(
                theme,
                'Failed to load subscription',
                error.toString(),
                () => ref.invalidate(subscriptionProvider),
              ),
              data: (subscription) {
                if (subscription == null) {
                  return _buildNoSubscription(context, theme);
                }
                return SubscriptionStatus(
                  subscription: subscription,
                  onChangePlan: () => _navigateToPlans(context),
                  onCancel: () => _showCancelConfirmation(context, ref),
                  onResume: () => _resumeSubscription(context, ref),
                );
              },
            ),
            const SizedBox(height: 24),

            // Invoices section
            Padding(
              padding: const EdgeInsets.only(left: 4, bottom: 8),
              child: Text(
                'Invoice History',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.onSurface.withOpacity(0.5),
                  letterSpacing: 0.5,
                ),
              ),
            ),
            invoicesAsync.when(
              loading: () => const Center(
                child: Padding(
                  padding: EdgeInsets.all(32),
                  child: CircularProgressIndicator(),
                ),
              ),
              error: (error, _) => _buildErrorCard(
                theme,
                'Failed to load invoices',
                error.toString(),
                () => ref.invalidate(invoicesProvider),
              ),
              data: (invoices) {
                if (invoices.isEmpty) {
                  return _buildEmptyInvoices(theme);
                }
                return _buildInvoiceList(theme, invoices);
              },
            ),
            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }

  Widget _buildNoSubscription(BuildContext context, ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          Icon(
            Icons.credit_card_off_outlined,
            size: 48,
            color: theme.colorScheme.onSurface.withOpacity(0.3),
          ),
          const SizedBox(height: 16),
          Text(
            'No Active Subscription',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Choose a plan to unlock premium features.',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.onSurface.withOpacity(0.6),
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 20),
          FilledButton.icon(
            onPressed: () => _navigateToPlans(context),
            icon: const Icon(Icons.rocket_launch_outlined),
            label: const Text('View Plans'),
            style: FilledButton.styleFrom(
              padding: const EdgeInsets.symmetric(
                horizontal: 24,
                vertical: 12,
              ),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyInvoices(ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          Icon(
            Icons.receipt_long_outlined,
            size: 40,
            color: theme.colorScheme.onSurface.withOpacity(0.3),
          ),
          const SizedBox(height: 12),
          Text(
            'No invoices yet',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.onSurface.withOpacity(0.6),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInvoiceList(ThemeData theme, List<Invoice> invoices) {
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
      ),
      clipBehavior: Clip.antiAlias,
      child: ListView.separated(
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        itemCount: invoices.length,
        separatorBuilder: (_, __) => Divider(
          height: 1,
          indent: 70,
          color: theme.colorScheme.outlineVariant,
        ),
        itemBuilder: (context, index) {
          return InvoiceTile(
            invoice: invoices[index],
            onDownload: invoices[index].downloadUrl != null
                ? () {
                    // TODO: Open download URL in browser
                  }
                : null,
          );
        },
      ),
    );
  }

  Widget _buildErrorCard(
    ThemeData theme,
    String title,
    String message,
    VoidCallback onRetry,
  ) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: theme.colorScheme.errorContainer.withOpacity(0.3),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          Icon(Icons.error_outline, color: theme.colorScheme.error),
          const SizedBox(height: 8),
          Text(title, style: theme.textTheme.titleSmall),
          const SizedBox(height: 4),
          Text(
            message,
            style: theme.textTheme.bodySmall?.copyWith(
              color: theme.colorScheme.onSurface.withOpacity(0.6),
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 12),
          TextButton.icon(
            onPressed: onRetry,
            icon: const Icon(Icons.refresh, size: 18),
            label: const Text('Retry'),
          ),
        ],
      ),
    );
  }

  void _navigateToPlans(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) => const PlansScreen()),
    );
  }

  Future<void> _showCancelConfirmation(
      BuildContext context, WidgetRef ref) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) {
        final theme = Theme.of(context);
        return AlertDialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          title: const Text('Cancel Subscription'),
          content: const Text(
            'Your subscription will remain active until the end of the '
            'current billing period. You can resume anytime before then.',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text('Keep Plan'),
            ),
            FilledButton(
              onPressed: () => Navigator.pop(context, true),
              style: FilledButton.styleFrom(
                backgroundColor: theme.colorScheme.error,
                foregroundColor: theme.colorScheme.onError,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: const Text('Cancel Plan'),
            ),
          ],
        );
      },
    );

    if (confirmed == true) {
      try {
        await ref.read(subscriptionProvider.notifier).cancel();
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: const Text('Subscription canceled'),
              behavior: SnackBarBehavior.floating,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          );
        }
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Failed to cancel: $e'),
              behavior: SnackBarBehavior.floating,
              backgroundColor: Theme.of(context).colorScheme.error,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          );
        }
      }
    }
  }

  Future<void> _resumeSubscription(
      BuildContext context, WidgetRef ref) async {
    try {
      await ref.read(subscriptionProvider.notifier).resume();
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: const Text('Subscription resumed'),
            behavior: SnackBarBehavior.floating,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to resume: $e'),
            behavior: SnackBarBehavior.floating,
            backgroundColor: Theme.of(context).colorScheme.error,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        );
      }
    }
  }
}

/// Skeleton placeholder while subscription data is loading.
class _SubscriptionSkeleton extends StatelessWidget {
  const _SubscriptionSkeleton();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
      ),
      child: const Center(
        child: Padding(
          padding: EdgeInsets.all(24),
          child: CircularProgressIndicator(),
        ),
      ),
    );
  }
}
