import 'package:flutter/material.dart';
import 'package:cleansaas_mobile/features/billing/models/plan.dart';

/// A card displaying a billing plan with pricing, features, and a CTA button.
///
/// The [isPopular] plan gets highlighted with a primary-colored border
/// and a "Popular" badge. The [isCurrentPlan] flag disables the CTA button
/// and shows "Current Plan" instead.
class PlanCard extends StatelessWidget {
  final Plan plan;
  final bool isCurrentPlan;
  final VoidCallback? onSelect;

  const PlanCard({
    super.key,
    required this.plan,
    this.isCurrentPlan = false,
    this.onSelect,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: plan.isPopular
              ? theme.colorScheme.primary
              : theme.colorScheme.outlineVariant,
          width: plan.isPopular ? 2 : 1,
        ),
      ),
      clipBehavior: Clip.antiAlias,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Popular badge
          if (plan.isPopular)
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(vertical: 6),
              color: theme.colorScheme.primary,
              child: Text(
                'Most Popular',
                textAlign: TextAlign.center,
                style: theme.textTheme.labelSmall?.copyWith(
                  color: theme.colorScheme.onPrimary,
                  fontWeight: FontWeight.bold,
                  letterSpacing: 0.5,
                ),
              ),
            ),

          Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Plan name
                Text(
                  plan.name,
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),

                // Price
                Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      plan.formattedPrice,
                      style: theme.textTheme.headlineLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: theme.colorScheme.primary,
                      ),
                    ),
                    if (plan.intervalLabel.isNotEmpty)
                      Padding(
                        padding: const EdgeInsets.only(bottom: 4, left: 2),
                        child: Text(
                          plan.intervalLabel,
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color:
                                theme.colorScheme.onSurface.withOpacity(0.5),
                          ),
                        ),
                      ),
                  ],
                ),
                const SizedBox(height: 20),

                // Features list
                ...plan.features.map(
                  (feature) => Padding(
                    padding: const EdgeInsets.only(bottom: 10),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Icon(
                          Icons.check_circle,
                          size: 18,
                          color: theme.colorScheme.primary,
                        ),
                        const SizedBox(width: 10),
                        Expanded(
                          child: Text(
                            feature,
                            style: theme.textTheme.bodyMedium,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                // CTA button
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: isCurrentPlan
                      ? OutlinedButton(
                          onPressed: null,
                          style: OutlinedButton.styleFrom(
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                          child: const Text('Current Plan'),
                        )
                      : plan.isPopular
                          ? FilledButton(
                              onPressed: onSelect,
                              style: FilledButton.styleFrom(
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              child: const Text(
                                'Get Started',
                                style: TextStyle(fontWeight: FontWeight.w600),
                              ),
                            )
                          : OutlinedButton(
                              onPressed: onSelect,
                              style: OutlinedButton.styleFrom(
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              child: const Text(
                                'Select Plan',
                                style: TextStyle(fontWeight: FontWeight.w600),
                              ),
                            ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
