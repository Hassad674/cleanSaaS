"use client";

import Link from "next/link";
import type { Notification } from "@/features/notification/types";

type NotificationDropdownProps = {
  notifications: Notification[];
  loading: boolean;
  onMarkRead: (id: string) => Promise<{ error: string | null }>;
  onMarkAllRead: () => Promise<{ error: string | null }>;
  onClose: () => void;
};

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
  const weeks = Math.floor(days / 7);
  if (weeks < 4) return `${weeks}w ago`;
  const months = Math.floor(days / 30);
  return `${months}mo ago`;
}

function LoadingSkeleton() {
  return (
    <div className="p-3 space-y-3">
      {[1, 2, 3].map((i) => (
        <div key={i} className="animate-pulse space-y-2">
          <div className="h-3 bg-muted rounded w-3/4" />
          <div className="h-2.5 bg-muted rounded w-full" />
        </div>
      ))}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="p-6 text-center">
      <svg
        className="mx-auto h-8 w-8 text-muted-foreground mb-2"
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
      <p className="text-sm text-muted-foreground">No notifications yet</p>
    </div>
  );
}

export function NotificationDropdown({
  notifications,
  loading,
  onMarkRead,
  onMarkAllRead,
  onClose,
}: NotificationDropdownProps) {
  const hasUnread = notifications.some((n) => !n.read);

  return (
    <div className="absolute right-0 top-full mt-2 w-80 sm:w-96 bg-card border border-border rounded-xl shadow-lg z-50 overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-border">
        <h3 className="text-sm font-medium text-foreground">Notifications</h3>
        {hasUnread && (
          <button
            onClick={() => onMarkAllRead()}
            className="text-xs text-primary hover:underline transition-colors"
          >
            Mark all as read
          </button>
        )}
      </div>

      {/* Content */}
      {loading ? (
        <LoadingSkeleton />
      ) : notifications.length === 0 ? (
        <EmptyState />
      ) : (
        <div className="max-h-80 overflow-y-auto">
          {notifications.map((notification) => (
            <button
              key={notification.id}
              onClick={() => {
                if (!notification.read) {
                  onMarkRead(notification.id);
                }
              }}
              className="w-full text-left px-4 py-3 hover:bg-muted/50 transition-colors border-b border-border last:border-b-0"
            >
              <div className="flex items-start gap-3">
                {/* Unread indicator */}
                <div className="mt-1.5 flex-shrink-0">
                  {!notification.read ? (
                    <span className="block h-2 w-2 rounded-full bg-primary" />
                  ) : (
                    <span className="block h-2 w-2" />
                  )}
                </div>

                <div className="flex-1 min-w-0">
                  <p
                    className={`text-sm truncate ${
                      notification.read
                        ? "text-muted-foreground"
                        : "text-foreground font-medium"
                    }`}
                  >
                    {notification.title}
                  </p>
                  <p className="text-xs text-muted-foreground mt-0.5 line-clamp-2">
                    {notification.message}
                  </p>
                  <p className="text-[11px] text-muted-foreground mt-1">
                    {timeAgo(notification.created_at)}
                  </p>
                </div>
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Footer */}
      <div className="border-t border-border px-4 py-2.5">
        <Link
          href="/notifications"
          onClick={onClose}
          className="text-xs text-primary hover:underline transition-colors"
        >
          View all notifications
        </Link>
      </div>
    </div>
  );
}
