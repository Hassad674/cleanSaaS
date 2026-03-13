"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type NotificationType = "system" | "billing" | "security" | "storage" | "ai";

interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message: string;
  read: boolean;
  createdAt: Date;
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

const INITIAL_NOTIFICATIONS: Notification[] = [
  {
    id: "n-1",
    type: "system",
    title: "Welcome to CleanSaaS!",
    message:
      "Your account has been set up. Explore features, customize your settings, and start building.",
    read: true,
    createdAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
  },
  {
    id: "n-2",
    type: "billing",
    title: "Your Pro plan is now active",
    message:
      "You have been upgraded to the Pro plan. Enjoy unlimited projects and priority support.",
    read: true,
    createdAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
  },
  {
    id: "n-3",
    type: "security",
    title: "New login detected from Chrome on macOS",
    message:
      "A new sign-in to your account was detected. If this was not you, change your password immediately.",
    read: false,
    createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000),
  },
  {
    id: "n-4",
    type: "storage",
    title: "File upload complete: report.pdf",
    message:
      "Your file report.pdf (2.4 MB) has been uploaded successfully and is ready to share.",
    read: false,
    createdAt: new Date(Date.now() - 3 * 60 * 60 * 1000),
  },
  {
    id: "n-5",
    type: "ai",
    title: "Your conversation 'Code Review' has a new response",
    message:
      "The AI assistant has finished analyzing your code and generated suggestions.",
    read: false,
    createdAt: new Date(Date.now() - 4 * 60 * 60 * 1000),
  },
  {
    id: "n-6",
    type: "system",
    title: "Scheduled maintenance on March 15",
    message:
      "We will perform a brief maintenance window on March 15 from 2:00 AM to 4:00 AM UTC.",
    read: false,
    createdAt: new Date(Date.now() - 6 * 60 * 60 * 1000),
  },
  {
    id: "n-7",
    type: "billing",
    title: "Invoice #1234 is ready",
    message:
      "Your invoice for February has been generated. View or download it from your billing settings.",
    read: true,
    createdAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
  },
  {
    id: "n-8",
    type: "security",
    title: "Password changed successfully",
    message:
      "Your password was updated. If you did not make this change, contact support immediately.",
    read: true,
    createdAt: new Date(Date.now() - 4 * 24 * 60 * 60 * 1000),
  },
  {
    id: "n-9",
    type: "system",
    title: "New feature: Dark mode is now available",
    message:
      "Switch between light and dark themes in your settings. Your preference is saved automatically.",
    read: true,
    createdAt: new Date(Date.now() - 6 * 24 * 60 * 60 * 1000),
  },
  {
    id: "n-10",
    type: "ai",
    title: "Usage limit: 80% of AI credits used",
    message:
      "You have used 80% of your monthly AI credits. Consider upgrading your plan for more capacity.",
    read: false,
    createdAt: new Date(Date.now() - 1 * 60 * 60 * 1000),
  },
];

const INCOMING_NOTIFICATIONS: Pick<Notification, "type" | "title" | "message">[] = [
  {
    type: "storage",
    title: "New file shared with you",
    message: "Alex shared 'Q1-Report.xlsx' with you. Open it from your files.",
  },
  {
    type: "system",
    title: "System update complete",
    message: "CleanSaaS v2.4 has been deployed. Check the changelog for details.",
  },
  {
    type: "ai",
    title: "AI model upgrade available",
    message:
      "GPT-4o is now available in your workspace. Switch models in AI settings.",
  },
  {
    type: "billing",
    title: "Payment method expiring soon",
    message:
      "Your Visa ending in 4242 expires next month. Update your card to avoid interruption.",
  },
  {
    type: "security",
    title: "Two-factor authentication reminder",
    message:
      "Protect your account by enabling 2FA. Set it up in your security settings.",
  },
  {
    type: "storage",
    title: "Storage quota at 90%",
    message:
      "You are using 4.5 GB of your 5 GB storage limit. Delete unused files or upgrade.",
  },
  {
    type: "system",
    title: "Invitation accepted",
    message: "Jordan accepted your team invitation and joined the workspace.",
  },
  {
    type: "ai",
    title: "Conversation export ready",
    message:
      "Your exported conversation 'Architecture Review' is ready to download.",
  },
];

const PAGE_SIZE = 5;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function relativeTime(date: Date): string {
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

const TYPE_LABELS: Record<NotificationType, string> = {
  system: "System",
  billing: "Billing",
  security: "Security",
  storage: "Storage",
  ai: "AI",
};

// ---------------------------------------------------------------------------
// Icons
// ---------------------------------------------------------------------------

function BellIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
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
  );
}

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M4.5 12.75l6 6 9-13.5"
      />
    </svg>
  );
}

function ChevronLeftIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M15.75 19.5L8.25 12l7.5-7.5"
      />
    </svg>
  );
}

function ChevronRightIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M8.25 4.5l7.5 7.5-7.5 7.5"
      />
    </svg>
  );
}

function ArrowLeftIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
      />
    </svg>
  );
}

function TypeIcon({ type, className }: { type: NotificationType; className?: string }) {
  switch (type) {
    case "system":
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 010 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 010-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28z" />
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      );
    case "billing":
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z" />
        </svg>
      );
    case "security":
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
        </svg>
      );
    case "storage":
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
        </svg>
      );
    case "ai":
      return (
        <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z" />
        </svg>
      );
  }
}

// ---------------------------------------------------------------------------
// Component
// ---------------------------------------------------------------------------

export function NotificationsDemo() {
  const [notifications, setNotifications] = useState<Notification[]>(INITIAL_NOTIFICATIONS);
  const [filter, setFilter] = useState<"all" | "unread">("all");
  const [page, setPage] = useState(0);
  const [newIds, setNewIds] = useState<Set<string>>(new Set());
  const incomingIndexRef = useRef(0);

  // ---- Derived state ----
  const unreadCount = useMemo(
    () => notifications.filter((n) => !n.read).length,
    [notifications],
  );

  const filtered = useMemo(
    () =>
      filter === "unread"
        ? notifications.filter((n) => !n.read)
        : notifications,
    [notifications, filter],
  );

  const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE));
  const paginated = filtered.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE);

  // Reset to first page when filter changes
  useEffect(() => {
    setPage(0);
  }, [filter]);

  // Clamp page if it goes out of bounds (e.g. after marking all read while on unread filter)
  useEffect(() => {
    if (page >= totalPages) {
      setPage(Math.max(0, totalPages - 1));
    }
  }, [page, totalPages]);

  // ---- Real-time simulation ----
  useEffect(() => {
    const interval = setInterval(() => {
      const idx = incomingIndexRef.current % INCOMING_NOTIFICATIONS.length;
      const template = INCOMING_NOTIFICATIONS[idx];
      const newNotification: Notification = {
        id: `n-incoming-${Date.now()}`,
        type: template.type,
        title: template.title,
        message: template.message,
        read: false,
        createdAt: new Date(),
      };

      incomingIndexRef.current += 1;

      setNotifications((current) => [newNotification, ...current]);
      setNewIds((current) => new Set(current).add(newNotification.id));

      // Remove from "new" set after animation completes
      setTimeout(() => {
        setNewIds((current) => {
          const next = new Set(current);
          next.delete(newNotification.id);
          return next;
        });
      }, 600);
    }, 15000);

    return () => clearInterval(interval);
  }, []);

  // ---- Actions ----
  const markAsRead = useCallback((id: string) => {
    setNotifications((current) =>
      current.map((n) => (n.id === id ? { ...n, read: true } : n)),
    );
  }, []);

  const markAllAsRead = useCallback(() => {
    setNotifications((current) => current.map((n) => ({ ...n, read: true })));
  }, []);

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 max-w-2xl">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeftIcon className="h-4 w-4" />
        Back to demos
      </Link>

      {/* Notification center card */}
      <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border">
          <div className="flex items-center gap-3">
            <div className="relative">
              <BellIcon className="h-6 w-6 text-foreground" />
              {unreadCount > 0 && (
                <span className="absolute -top-1.5 -right-1.5 flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1 text-[11px] font-semibold text-primary-foreground">
                  {unreadCount}
                </span>
              )}
            </div>
            <h1 className="text-lg font-semibold text-foreground">
              Notifications
            </h1>
          </div>

          {unreadCount > 0 && (
            <button
              type="button"
              onClick={markAllAsRead}
              className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
            >
              <CheckIcon className="h-4 w-4" />
              Mark all read
            </button>
          )}
        </div>

        {/* Filter tabs */}
        <div className="flex border-b border-border">
          <button
            type="button"
            onClick={() => setFilter("all")}
            className={cn(
              "flex-1 py-2.5 text-sm font-medium text-center transition-colors",
              filter === "all"
                ? "text-primary border-b-2 border-primary"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            All
          </button>
          <button
            type="button"
            onClick={() => setFilter("unread")}
            className={cn(
              "flex-1 py-2.5 text-sm font-medium text-center transition-colors",
              filter === "unread"
                ? "text-primary border-b-2 border-primary"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            Unread
            {unreadCount > 0 && (
              <span className="ml-1.5 inline-flex h-5 min-w-5 items-center justify-center rounded-full bg-primary/10 px-1.5 text-xs font-semibold text-primary">
                {unreadCount}
              </span>
            )}
          </button>
        </div>

        {/* Notification list */}
        <div className="divide-y divide-border">
          {paginated.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 px-4">
              <BellIcon className="h-10 w-10 text-muted-foreground/50 mb-3" />
              <p className="text-sm text-muted-foreground">
                {filter === "unread"
                  ? "No unread notifications"
                  : "No notifications yet"}
              </p>
            </div>
          ) : (
            paginated.map((notification) => (
              <NotificationRow
                key={notification.id}
                notification={notification}
                isNew={newIds.has(notification.id)}
                onMarkRead={markAsRead}
              />
            ))
          )}
        </div>

        {/* Pagination */}
        {filtered.length > PAGE_SIZE && (
          <div className="flex items-center justify-between px-5 py-3 border-t border-border bg-muted/30">
            <button
              type="button"
              onClick={() => setPage((p) => Math.max(0, p - 1))}
              disabled={page === 0}
              className={cn(
                "inline-flex items-center gap-1 text-sm font-medium transition-colors",
                page === 0
                  ? "text-muted-foreground/40 cursor-not-allowed"
                  : "text-foreground hover:text-primary",
              )}
            >
              <ChevronLeftIcon className="h-4 w-4" />
              Previous
            </button>

            <span className="text-xs text-muted-foreground">
              Page {page + 1} of {totalPages}
            </span>

            <button
              type="button"
              onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
              disabled={page >= totalPages - 1}
              className={cn(
                "inline-flex items-center gap-1 text-sm font-medium transition-colors",
                page >= totalPages - 1
                  ? "text-muted-foreground/40 cursor-not-allowed"
                  : "text-foreground hover:text-primary",
              )}
            >
              Next
              <ChevronRightIcon className="h-4 w-4" />
            </button>
          </div>
        )}
      </div>

      {/* Footer note */}
      <p className="text-center text-sm text-muted-foreground mt-6">
        New notifications appear every 15 seconds. In production, these arrive
        via WebSocket from the Go backend.
      </p>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Notification row
// ---------------------------------------------------------------------------

function NotificationRow({
  notification,
  isNew,
  onMarkRead,
}: {
  notification: Notification;
  isNew: boolean;
  onMarkRead: (id: string) => void;
}) {
  return (
    <div
      role="button"
      tabIndex={0}
      onClick={() => {
        if (!notification.read) onMarkRead(notification.id);
      }}
      onKeyDown={(e) => {
        if ((e.key === "Enter" || e.key === " ") && !notification.read) {
          e.preventDefault();
          onMarkRead(notification.id);
        }
      }}
      className={cn(
        "flex gap-3.5 px-5 py-4 transition-all duration-200",
        !notification.read && "border-l-2 border-l-primary bg-primary/[0.03]",
        notification.read && "border-l-2 border-l-transparent",
        !notification.read && "cursor-pointer hover:bg-accent/50",
        notification.read && "cursor-default",
        isNew && "animate-[slideIn_0.4s_ease-out_forwards]",
      )}
    >
      {/* Icon + dot */}
      <div className="relative flex-shrink-0 mt-0.5">
        <div
          className={cn(
            "flex h-9 w-9 items-center justify-center rounded-lg",
            !notification.read ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground",
          )}
        >
          <TypeIcon type={notification.type} className="h-5 w-5" />
        </div>
        {!notification.read && (
          <span className="absolute -top-0.5 -right-0.5 h-2.5 w-2.5 rounded-full bg-primary ring-2 ring-card" />
        )}
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <p
              className={cn(
                "text-sm leading-snug",
                !notification.read
                  ? "font-semibold text-foreground"
                  : "font-medium text-foreground/80",
              )}
            >
              {notification.title}
            </p>
            <p className="text-sm text-muted-foreground mt-0.5 line-clamp-2 leading-relaxed">
              {notification.message}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3 mt-2">
          <span className="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-[11px] font-medium text-muted-foreground">
            {TYPE_LABELS[notification.type]}
          </span>
          <span className="text-xs text-muted-foreground">
            {relativeTime(notification.createdAt)}
          </span>
          {!notification.read && (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                onMarkRead(notification.id);
              }}
              className="ml-auto text-xs font-medium text-primary hover:text-primary/80 transition-colors"
            >
              Mark read
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
