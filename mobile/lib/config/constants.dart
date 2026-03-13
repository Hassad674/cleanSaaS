/// Application-wide constants.
///
/// Edit [apiBaseUrl] to point to your backend instance.
class AppConstants {
  AppConstants._();

  /// App identity
  static const String appName = 'CleanSaaS';
  static const String appDescription =
      'Open-source SaaS boilerplate for modern applications';

  /// API configuration
  ///
  /// Android emulator: http://10.0.2.2:8081
  /// iOS simulator:    http://localhost:8081
  /// Physical device:  http://<your-local-ip>:8081
  /// Production:       https://api.yourdomain.com
  static const String apiBaseUrl = 'http://10.0.2.2:8081';

  /// HTTP timeouts (milliseconds)
  static const int connectTimeout = 15000;
  static const int receiveTimeout = 15000;
  static const int sendTimeout = 15000;

  /// Secure storage keys
  static const String tokenKey = 'cleansaas_token';
  static const String userKey = 'cleansaas_user';

  /// Pagination
  static const int defaultPageSize = 20;
  static const int maxPageSize = 100;
}
