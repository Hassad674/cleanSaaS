/// Represents the user's current subscription.
///
/// Tracks the plan, status, and billing period dates.
class Subscription {
  final String id;
  final String planId;
  final String planName;
  final String status; // 'active', 'canceled', 'past_due', 'trialing'
  final DateTime currentPeriodStart;
  final DateTime currentPeriodEnd;
  final DateTime? canceledAt;
  final DateTime createdAt;

  const Subscription({
    required this.id,
    required this.planId,
    required this.planName,
    required this.status,
    required this.currentPeriodStart,
    required this.currentPeriodEnd,
    this.canceledAt,
    required this.createdAt,
  });

  factory Subscription.fromJson(Map<String, dynamic> json) {
    return Subscription(
      id: json['id'] as String,
      planId: json['plan_id'] as String,
      planName: json['plan_name'] as String,
      status: json['status'] as String,
      currentPeriodStart:
          DateTime.parse(json['current_period_start'] as String),
      currentPeriodEnd: DateTime.parse(json['current_period_end'] as String),
      canceledAt: json['canceled_at'] != null
          ? DateTime.parse(json['canceled_at'] as String)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'plan_id': planId,
      'plan_name': planName,
      'status': status,
      'current_period_start': currentPeriodStart.toIso8601String(),
      'current_period_end': currentPeriodEnd.toIso8601String(),
      'canceled_at': canceledAt?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
    };
  }

  /// Whether the subscription is currently active and not canceled.
  bool get isActive => status == 'active' || status == 'trialing';

  /// Whether the subscription has been canceled but is still in the paid period.
  bool get isCanceledButActive =>
      status == 'canceled' && currentPeriodEnd.isAfter(DateTime.now());

  /// Returns a human-readable status label.
  String get statusLabel {
    switch (status) {
      case 'active':
        return 'Active';
      case 'canceled':
        return isCanceledButActive ? 'Cancels soon' : 'Canceled';
      case 'past_due':
        return 'Past due';
      case 'trialing':
        return 'Trial';
      default:
        return status;
    }
  }
}

/// Represents an invoice/payment record.
class Invoice {
  final String id;
  final int amountCents;
  final String currency;
  final String status; // 'paid', 'pending', 'failed'
  final DateTime createdAt;
  final String? downloadUrl;

  const Invoice({
    required this.id,
    required this.amountCents,
    required this.currency,
    required this.status,
    required this.createdAt,
    this.downloadUrl,
  });

  factory Invoice.fromJson(Map<String, dynamic> json) {
    return Invoice(
      id: json['id'] as String,
      amountCents: json['amount_cents'] as int,
      currency: json['currency'] as String? ?? 'usd',
      status: json['status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      downloadUrl: json['download_url'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'amount_cents': amountCents,
      'currency': currency,
      'status': status,
      'created_at': createdAt.toIso8601String(),
      'download_url': downloadUrl,
    };
  }

  /// Formats the amount in dollars (e.g., "\$9.99").
  String get formattedAmount {
    final dollars = amountCents / 100;
    return '\$${dollars.toStringAsFixed(2)}';
  }

  /// Returns a human-readable status label.
  String get statusLabel {
    switch (status) {
      case 'paid':
        return 'Paid';
      case 'pending':
        return 'Pending';
      case 'failed':
        return 'Failed';
      default:
        return status;
    }
  }
}
