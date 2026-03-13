import type { Metadata } from "next";
import { ChatLayout } from "@/features/ai/components/chat-layout";

export const metadata: Metadata = { title: "AI Chat" };

export default function AIPage() {
  return <ChatLayout />;
}
