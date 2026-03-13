/// Represents a billing plan available for subscription.
///
/// Plans have a name, price, billing interval, and a list of included features.
class Plan {
  final String id;
  final String name;
  final int priceCents;
  final String interval; // 'month' or 'year'
  final List<String> features;
  final bool isPopular;

  const Plan({
    required this.id,
    required this.name,
    required this.priceCents,
    required this.interval,
    required this.features,
    this.isPopular = false,
  });

  factory Plan.fromJson(Map<String, dynamic> json) {
    return Plan(
      id: json['id'] as String,
      name: json['name'] as String,
      priceCents: json['price_cents'] as int,
      interval: json['interval'] as String,
      features: (json['features'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      isPopular: json['is_popular'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'price_cents': priceCents,
      'interval': interval,
      'features': features,
      'is_popular': isPopular,
    };
  }

  /// Formats the price in dollars (e.g., "\$9.99").
  String get formattedPrice {
    if (priceCents == 0) return 'Free';
    final dollars = priceCents / 100;
    return '\$${dollars.toStringAsFixed(dollars.truncateToDouble() == dollars ? 0 : 2)}';
  }

  /// Returns the billing period label (e.g., "/month" or "/year").
  String get intervalLabel {
    if (priceCents == 0) return '';
    return '/$interval';
  }
}
