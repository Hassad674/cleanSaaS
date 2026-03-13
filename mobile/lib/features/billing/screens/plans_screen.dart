import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:cleansaas_mobile/features/billing/providers/billing_provider.dart';
import 'package:cleansaas_mobile/features/billing/widgets/plan_card.dart';

/// Screen displaying all available billing plans in a scrollable list.
///
/// Highlights the current plan and allows users to select a new plan,
/// which initiates a checkout flow via an external browser.
class PlansScreen extends ConsumerWidget {
  const PlansScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final plansAsync = ref.watch(plansProvider);
    final subscriptionAsync = ref.watch(subscriptionProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Choose a Plan'),
      ),
      body: plansAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => _buildErrorState(context, ref, theme, error),
        data: (plans) {
          final currentPlanId = subscriptionAsync.valueOrNull?.planId;

          return LayoutBuilder(
            builder: (context, constraints) {
              // Tablet layout: side-by-side cards.
              if (constraints.maxWidth > 700) {
                return SingleChildScrollView(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      _buildHeader(theme),
                      const SizedBox(height: 24),
                      Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: plans.map((plan) {
                          return Expanded(
                            child: Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 6),
                              child: PlanCard(
                                plan: plan,
                                isCurrentPlan: plan.id == currentPlanId,
                                onSelect: plan.id == currentPlanId
                                    ? null
                                    : () => _selectPlan(
                                          context,
                                          ref,
                                          plan.id,
                                        ),
                              ),
                            ),
                          );
                        }).toList(),
                      ),
                    ],
                  ),
                );
              }

              // Phone layout: vertical list.
              return ListView.builder(
                padding: const EdgeInsets.all(16),
                itemCount: plans.length + 1, // +1 for header
                itemBuilder: (context, index) {
                  if (index == 0) {
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 20),
                      child: _buildHeader(theme),
                    );
                  }

                  final plan = plans[index - 1];
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 16),
                    child: PlanCard(
                      plan: plan,
                      isCurrentPlan: plan.id == currentPlanId,
                      onSelect: plan.id == currentPlanId
                          ? null
                          : () => _selectPlan(context, ref, plan.id),
                    ),
                  );
                },
              );
            },
          );
        },
      ),
    );
  }

  Widget _buildHeader(ThemeData theme) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Find the right plan for you',
          style: theme.textTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          'Upgrade or downgrade anytime. No hidden fees.',
          style: theme.textTheme.bodyMedium?.copyWith(
            color: theme.colorScheme.onSurface.withOpacity(0.6),
          ),
        ),
      ],
    );
  }

  Widget _buildErrorState(
    BuildContext context,
    WidgetRef ref,
    ThemeData theme,
    Object error,
  ) {
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
              'Failed to load plans',
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
              onPressed: () => ref.invalidate(plansProvider),
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _selectPlan(
    BuildContext context,
    WidgetRef ref,
    String planId,
  ) async {
    try {
      final checkoutUrl = await ref
          .read(subscriptionProvider.notifier)
          .checkout(planId: planId);

      // Open checkout URL in external browser.
      final uri = Uri.parse(checkoutUrl);
      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      } else {
        throw Exception('Could not open checkout page');
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Checkout failed: $e'),
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
