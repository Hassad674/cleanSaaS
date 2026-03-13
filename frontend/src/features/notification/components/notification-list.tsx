"use client";

import { useNotifications } from "@/features/notification/hooks/use-notifications";
import { formatDate } from "@/shared/lib/utils";
import type { Notification } from "@/features/notification/types";

function timeAgo(dateStr: string): string {
  const now = Date.now();
  const date = new Date(dateStr).getTime();
  const seconds = Math.floor((now - date) / 1000);

  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;
  return formatDate(dateStr);
}

function LoadingSkeleton() {
  return (
    <div className="space-y-3">
      {[1, 2, 3, 4, 5].map((i) => (
        <div
          key={i}
          className="bg-card border border-border rounded-xl p-4 shadow-sm animate-pulse"
        >
          <div className="flex items-start gap-3">
            <div className="h-2 w-2 rounded-full bg-muted mt-2" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-muted rounded w-1/3" />
              <div className="h-3 bg-muted rounded w-2/3" />
              <div className="h-2.5 bg-muted rounded w-1/4" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

function EmptyState({ unreadOnly }: { unreadOnly: boolean }) {
  return (
    <div className="bg-card border border-border rounded-xl p-12 shadow-sm text-center">
      <svg
        className="mx-auto h-12 w-12 text-muted-foreground mb-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0"
        />
      </svg>
      <h3 className="text-base font-medium text-foreground mb-1">
        {unreadOnly ? "No unread notifications" : "No notifications yet"}
      </h3>
      <p className="text-sm text-muted-foreground">
        {unreadOnly
          ? "You're all caught up! Switch to all notifications to see your history."
          : "When you receive notifications, they will appear here."}
      </p>
    </div>
  );
}

function NotificationRow({
  notification,
  onMarkRead,
}: {
  notification: Notification;
  onMarkRead: (id: string) => Promise<{ error: string | null }>;
}) {
  return (
    <div
      className={`bg-card border border-border rounded-xl p-4 shadow-sm transition-colors ${
        !notification.read ? "border-l-primary border-l-2" : ""
      }`}
    >
      <div className="flex items-start gap-3">
        {/* Unread indicator */}
        <div className="mt-1.5 flex-shrink-0">
          {!notification.read ? (
            <span className="block h-2.5 w-2.5 rounded-full bg-primary" />
          ) : (
            <span className="block h-2.5 w-2.5 rounded-full bg-muted" />
          )}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0">
              <p
                className={`text-sm ${
                  notification.read
                    ? "text-muted-foreground"
                    : "text-foreground font-medium"
                }`}
              >
                {notification.title}
              </p>
              <p className="text-sm text-muted-foreground mt-1">
                {notification.message}
              </p>
              <p className="text-xs text-muted-foreground mt-2">
                {timeAgo(notification.created_at)}
              </p>
            </div>

            {/* Mark as read button */}
            {!notification.read && (
              <button
                onClick={() => onMarkRead(notification.id)}
                className="flex-shrink-0 text-xs text-primary hover:underline transition-colors"
              >
                Mark read
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export function NotificationList() {
  const {
    notifications,
    total,
    page,
    totalPages,
    unreadCount,
    unreadOnly,
    loading,
    error,
    hasNext,
    hasPrev,
    markRead,
    markAllRead,
    toggleUnreadOnly,
    goToNextPage,
    goToPrevPage,
  } = useNotifications();

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Notifications</h1>
          <p className="text-muted-foreground mt-1">
            {unreadCount > 0
              ? `You have ${unreadCount} unread notification${unreadCount !== 1 ? "s" : ""}.`
              : "You're all caught up."}
          </p>
        </div>

        <div className="flex items-center gap-3">
          {/* Filter toggle */}
          <button
            onClick={toggleUnreadOnly}
            className={`text-sm px-3 py-1.5 rounded-lg transition-colors ${
              unreadOnly
                ? "bg-primary text-primary-foreground"
                : "bg-muted text-muted-foreground hover:text-foreground"
            }`}
          >
            {unreadOnly ? "Unread only" : "All"}
          </button>

          {/* Mark all as read */}
          {unreadCount > 0 && (
            <button
              onClick={() => markAllRead()}
              className="text-sm text-primary hover:underline transition-colors"
            >
              Mark all as read
            </button>
          )}
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3">
          <p className="text-sm text-destructive">{error}</p>
        </div>
      )}

      {/* Notification list */}
      {loading ? (
        <LoadingSkeleton />
      ) : notifications.length === 0 ? (
        <EmptyState unreadOnly={unreadOnly} />
      ) : (
        <div className="space-y-3">
          {notifications.map((notification) => (
            <NotificationRow
              key={notification.id}
              notification={notification}
              onMarkRead={markRead}
            />
          ))}
        </div>
      )}

      {/* Pagination */}
      {(hasPrev || hasNext) && (
        <div className="flex items-center justify-between pt-2">
          <button
            onClick={goToPrevPage}
            disabled={!hasPrev}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
          >
            Previous
          </button>
          <span className="text-xs text-muted-foreground">
            Page {page} of {totalPages} &middot; {total} notification
            {total !== 1 ? "s" : ""}
          </span>
          <button
            onClick={goToNextPage}
            disabled={!hasNext}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
