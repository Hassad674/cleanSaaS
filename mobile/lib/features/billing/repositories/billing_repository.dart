import 'package:cleansaas_mobile/core/api/api_client.dart';
import 'package:cleansaas_mobile/features/billing/models/plan.dart';
import 'package:cleansaas_mobile/features/billing/models/subscription.dart';

/// Repository handling all billing-related API operations.
///
/// Provides methods to fetch plans, manage subscriptions, list invoices,
/// initiate checkout, and cancel subscriptions.
class BillingRepository {
  final ApiClient _apiClient;

  BillingRepository(this._apiClient);

  /// Fetches all available billing plans.
  Future<List<Plan>> getPlans() async {
    final response = await _apiClient.get('/billing/plans');
    final plans = (response.data as List<dynamic>)
        .map((json) => Plan.fromJson(json as Map<String, dynamic>))
        .toList();
    return plans;
  }

  /// Fetches the current user's active subscription.
  ///
  /// Returns null if the user has no active subscription.
  Future<Subscription?> getSubscription() async {
    try {
      final response = await _apiClient.get('/billing/subscription');
      if (response.data == null) return null;
      return Subscription.fromJson(response.data as Map<String, dynamic>);
    } catch (_) {
      // No active subscription.
      return null;
    }
  }

  /// Fetches the user's invoice history.
  Future<List<Invoice>> getInvoices() async {
    final response = await _apiClient.get('/billing/invoices');
    final invoices = (response.data as List<dynamic>)
        .map((json) => Invoice.fromJson(json as Map<String, dynamic>))
        .toList();
    return invoices;
  }

  /// Initiates a checkout session for the given plan.
  ///
  /// Returns the checkout URL that should be opened in a browser.
  Future<String> checkout({required String planId}) async {
    final response = await _apiClient.post('/billing/checkout', data: {
      'plan_id': planId,
    });
    return response.data['checkout_url'] as String;
  }

  /// Cancels the current subscription.
  ///
  /// The subscription remains active until the end of the current billing period.
  Future<void> cancelSubscription() async {
    await _apiClient.post('/billing/cancel');
  }

  /// Resumes a previously canceled subscription.
  Future<void> resumeSubscription() async {
    await _apiClient.post('/billing/resume');
  }
}
