import type { Metadata } from "next";
import { AdminDemo } from "./admin-demo";

export const metadata: Metadata = {
  title: "Admin Demo — CleanSaaS",
  description:
    "Interactive demo of CleanSaaS admin dashboard: user management, analytics, and system monitoring. No account required.",
};

export default function AdminDemoPage() {
  return <AdminDemo />;
}
