"use server";

import { api } from "@/shared/lib/api";
import type {
  NotificationsResponse,
  UnreadCountResponse,
} from "@/features/notification/types";

export async function getNotifications(
  authToken: string,
  page: number = 1,
  limit: number = 20,
  unreadOnly: boolean = false
) {
  const params = new URLSearchParams({
    page: String(page),
    limit: String(limit),
  });
  if (unreadOnly) {
    params.set("unread", "true");
  }

  return api<NotificationsResponse>(`/notifications?${params.toString()}`, {
    token: authToken,
  });
}

export async function getUnreadCount(authToken: string) {
  return api<UnreadCountResponse>("/notifications/count", {
    token: authToken,
  });
}

export async function markAsRead(notifId: string, authToken: string) {
  return api<null>(`/notifications/${notifId}/read`, {
    method: "PUT",
    token: authToken,
  });
}

export async function markAllAsRead(authToken: string) {
  return api<null>("/notifications/read-all", {
    method: "PUT",
    token: authToken,
  });
}
