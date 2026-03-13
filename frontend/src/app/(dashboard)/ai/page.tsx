import type { Metadata } from "next";
import { ChatLayout } from "@/features/ai/components/chat-layout";

export const metadata: Metadata = {
  title: "AI Chat",
  robots: { index: false, follow: false },
};

export default function AIPage() {
  return <ChatLayout />;
}
