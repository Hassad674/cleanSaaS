import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/auth/providers/auth_provider.dart';
import 'package:cleansaas_mobile/features/billing/models/plan.dart';
import 'package:cleansaas_mobile/features/billing/models/subscription.dart';
import 'package:cleansaas_mobile/features/billing/repositories/billing_repository.dart';

/// Provider for the [BillingRepository] instance.
final billingRepositoryProvider = Provider<BillingRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return BillingRepository(apiClient);
});

/// Async provider that fetches and caches available billing plans.
final plansProvider = FutureProvider<List<Plan>>((ref) async {
  final repository = ref.watch(billingRepositoryProvider);
  return repository.getPlans();
});

/// Async provider that fetches and caches the current subscription.
///
/// Invalidate this provider to trigger a fresh fetch after changes.
final subscriptionProvider =
    AsyncNotifierProvider<SubscriptionNotifier, Subscription?>(
  SubscriptionNotifier.new,
);

/// Notifier managing the subscription state.
class SubscriptionNotifier extends AsyncNotifier<Subscription?> {
  @override
  Future<Subscription?> build() async {
    final repository = ref.watch(billingRepositoryProvider);
    return repository.getSubscription();
  }

  /// Initiates checkout for a plan. Returns the checkout URL.
  Future<String> checkout({required String planId}) async {
    final repository = ref.read(billingRepositoryProvider);
    return repository.checkout(planId: planId);
  }

  /// Cancels the current subscription.
  Future<void> cancel() async {
    final repository = ref.read(billingRepositoryProvider);
    await repository.cancelSubscription();
    ref.invalidateSelf();
  }

  /// Resumes a canceled subscription.
  Future<void> resume() async {
    final repository = ref.read(billingRepositoryProvider);
    await repository.resumeSubscription();
    ref.invalidateSelf();
  }

  /// Forces a refresh of the subscription from the API.
  Future<void> refresh() async {
    ref.invalidateSelf();
  }
}

/// Async provider that fetches and caches invoices.
final invoicesProvider = FutureProvider<List<Invoice>>((ref) async {
  final repository = ref.watch(billingRepositoryProvider);
  return repository.getInvoices();
});
