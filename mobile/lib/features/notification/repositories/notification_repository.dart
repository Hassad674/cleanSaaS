import 'package:cleansaas_mobile/core/api/api_client.dart';
import 'package:cleansaas_mobile/features/notification/models/notification.dart';

/// Repository handling all notification API communication.
///
/// Provides methods for listing notifications, reading unread counts,
/// and marking notifications as read (individually or all at once).
class NotificationRepository {
  final ApiClient _apiClient;

  NotificationRepository({required ApiClient apiClient})
      : _apiClient = apiClient;

  /// Fetches a paginated list of notifications.
  ///
  /// Set [unreadOnly] to true to filter for only unread notifications.
  Future<List<AppNotification>> getNotifications({
    int page = 1,
    int limit = 20,
    bool unreadOnly = false,
  }) async {
    final response = await _apiClient.get(
      '/notifications',
      queryParameters: {
        'page': page.toString(),
        'limit': limit.toString(),
        if (unreadOnly) 'unread': 'true',
      },
    );

    final List<dynamic> data = response.data['data'] as List<dynamic>;
    return data
        .map((json) => AppNotification.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Returns the count of unread notifications.
  Future<int> getUnreadCount() async {
    final response = await _apiClient.get('/notifications/count');
    return response.data['data']['count'] as int;
  }

  /// Marks a single notification as read by its [notificationId].
  Future<void> markAsRead(String notificationId) async {
    await _apiClient.put('/notifications/$notificationId/read');
  }

  /// Marks all notifications as read for the current user.
  Future<void> markAllAsRead() async {
    await _apiClient.put('/notifications/read-all');
  }
}
