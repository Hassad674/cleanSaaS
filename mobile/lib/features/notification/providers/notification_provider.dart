import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cleansaas_mobile/features/notification/models/notification.dart';
import 'package:cleansaas_mobile/features/notification/repositories/notification_repository.dart';

/// State management for the notifications feature.
///
/// Manages the notifications list, unread count, pagination,
/// filtering, and read/unread state transitions.
class NotificationProvider extends ChangeNotifier {
  final NotificationRepository _repository;

  NotificationProvider({required NotificationRepository repository})
      : _repository = repository;

  // -- Notifications list state --

  List<AppNotification> _notifications = [];
  List<AppNotification> get notifications => _notifications;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _error;
  String? get error => _error;

  int _currentPage = 1;
  bool _hasMore = true;
  bool get hasMore => _hasMore;

  // -- Unread state --

  int _unreadCount = 0;
  int get unreadCount => _unreadCount;

  bool _showUnreadOnly = false;
  bool get showUnreadOnly => _showUnreadOnly;

  /// Loads the first page of notifications, replacing any existing data.
  Future<void> loadNotifications() async {
    _isLoading = true;
    _error = null;
    _currentPage = 1;
    notifyListeners();

    try {
      _notifications = await _repository.getNotifications(
        page: 1,
        unreadOnly: _showUnreadOnly,
      );
      _hasMore = _notifications.length >= 20;
      // Also refresh unread count.
      _unreadCount = await _repository.getUnreadCount();
    } catch (e) {
      _error = 'Failed to load notifications. Please try again.';
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  /// Loads the next page of notifications and appends them.
  Future<void> loadMoreNotifications() async {
    if (_isLoading || !_hasMore) return;

    _isLoading = true;
    notifyListeners();

    try {
      _currentPage++;
      final more = await _repository.getNotifications(
        page: _currentPage,
        unreadOnly: _showUnreadOnly,
      );
      _notifications = [..._notifications, ...more];
      _hasMore = more.length >= 20;
    } catch (e) {
      _currentPage--;
      _error = 'Failed to load more notifications.';
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  /// Toggles between showing all notifications and unread-only.
  Future<void> toggleUnreadFilter() async {
    _showUnreadOnly = !_showUnreadOnly;
    await loadNotifications();
  }

  /// Marks a single notification as read.
  ///
  /// Updates both the local state and the server.
  Future<void> markAsRead(String notificationId) async {
    try {
      await _repository.markAsRead(notificationId);

      _notifications = _notifications.map((n) {
        if (n.id == notificationId && !n.read) {
          _unreadCount = (_unreadCount - 1).clamp(0, _unreadCount);
          return n.copyWith(read: true);
        }
        return n;
      }).toList();

      notifyListeners();
    } catch (e) {
      _error = 'Failed to mark notification as read.';
      notifyListeners();
    }
  }

  /// Marks all notifications as read.
  ///
  /// Updates both the local state and the server.
  Future<void> markAllAsRead() async {
    try {
      await _repository.markAllAsRead();

      _notifications = _notifications.map((n) => n.copyWith(read: true)).toList();
      _unreadCount = 0;

      notifyListeners();
    } catch (e) {
      _error = 'Failed to mark all as read.';
      notifyListeners();
    }
  }

  /// Refreshes just the unread count without reloading the full list.
  Future<void> refreshUnreadCount() async {
    try {
      _unreadCount = await _repository.getUnreadCount();
      notifyListeners();
    } catch (_) {
      // Silently fail — unread count is not critical.
    }
  }
}

/// Riverpod provider for [NotificationProvider].
///
/// Requires a [NotificationRepository] to be provided in the widget tree.
final notificationProvider = ChangeNotifierProvider<NotificationProvider>((ref) {
  throw UnimplementedError(
    'notificationProvider must be overridden with a valid NotificationRepository.',
  );
});
