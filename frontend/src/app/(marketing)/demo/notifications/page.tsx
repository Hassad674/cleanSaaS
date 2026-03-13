import type { Metadata } from "next";
import { NotificationsDemo } from "./notifications-demo";

export const metadata: Metadata = {
  title: "Notifications Demo — CleanSaaS",
  description:
    "Interactive notification center demo with read/unread management, real-time simulation, and pagination. No account required.",
};

export default function NotificationsDemoPage() {
  return <NotificationsDemo />;
}
