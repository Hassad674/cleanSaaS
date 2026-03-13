"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { useAuth } from "@/features/auth/hooks/use-auth";
import {
  getNotifications as getNotificationsAction,
  getUnreadCount as getUnreadCountAction,
  markAsRead as markAsReadAction,
  markAllAsRead as markAllAsReadAction,
} from "@/features/notification/actions/notification";
import { PAGINATION_DEFAULT_LIMIT } from "@/shared/lib/constants";
import type { Notification } from "@/features/notification/types";

const POLL_INTERVAL_MS = 30_000;

export function useNotifications() {
  const { getToken } = useAuth({ required: true });

  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [unreadCount, setUnreadCount] = useState(0);
  const [unreadOnly, setUnreadOnly] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const limit = PAGINATION_DEFAULT_LIMIT;
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const fetchUnreadCount = useCallback(() => {
    const token = getToken();
    if (!token) return;

    getUnreadCountAction(token).then((res) => {
      if (res.data) {
        setUnreadCount(res.data.count);
      }
    });
  }, [getToken]);

  const fetchNotifications = useCallback(() => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);
    getNotificationsAction(token, page, limit, unreadOnly).then((res) => {
      if (res.data) {
        setNotifications(res.data.notifications ?? []);
        setTotal(res.data.total);
      } else {
        setError(res.error ?? "Failed to load notifications");
      }
      setLoading(false);
    });
  }, [getToken, page, limit, unreadOnly]);

  // Fetch notifications when page or filter changes
  useEffect(() => {
    fetchNotifications();
  }, [fetchNotifications]);

  // Fetch unread count on mount
  useEffect(() => {
    fetchUnreadCount();
  }, [fetchUnreadCount]);

  // Poll unread count every 30 seconds
  useEffect(() => {
    intervalRef.current = setInterval(fetchUnreadCount, POLL_INTERVAL_MS);
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [fetchUnreadCount]);

  const markRead = useCallback(
    async (notifId: string) => {
      const token = getToken();
      if (!token) return { error: "Not authenticated" };

      const res = await markAsReadAction(notifId, token);

      if (res.error) {
        setError(res.error);
        return { error: res.error };
      }

      // Optimistically update the notification in the list
      setNotifications((prev) =>
        prev.map((n) => (n.id === notifId ? { ...n, read: true } : n))
      );
      setUnreadCount((prev) => Math.max(0, prev - 1));

      return { error: null };
    },
    [getToken]
  );

  const markAllRead = useCallback(async () => {
    const token = getToken();
    if (!token) return { error: "Not authenticated" };

    const res = await markAllAsReadAction(token);

    if (res.error) {
      setError(res.error);
      return { error: res.error };
    }

    // Optimistically mark all as read
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    setUnreadCount(0);

    return { error: null };
  }, [getToken]);

  const toggleUnreadOnly = useCallback(() => {
    setUnreadOnly((prev) => !prev);
    setPage(1);
  }, []);

  const totalPages = Math.ceil(total / limit);
  const hasNext = page < totalPages;
  const hasPrev = page > 1;

  const goToNextPage = useCallback(() => {
    if (hasNext) setPage((prev) => prev + 1);
  }, [hasNext]);

  const goToPrevPage = useCallback(() => {
    if (hasPrev) setPage((prev) => prev - 1);
  }, [hasPrev]);

  return {
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
    fetchNotifications,
    fetchUnreadCount,
  };
}
