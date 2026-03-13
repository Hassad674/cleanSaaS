import type { Metadata } from "next";
import { NotificationList } from "@/features/notification/components/notification-list";

export const metadata: Metadata = { title: "Notifications" };

export default function NotificationsPage() {
  return <NotificationList />;
}
