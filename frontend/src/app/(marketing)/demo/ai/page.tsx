import type { Metadata } from "next";
import { AiDemo } from "./ai-demo";

export const metadata: Metadata = {
  title: "AI Chat Demo — CleanSaaS",
  description:
    "Interactive AI chat demo with conversation history, streaming responses, and multi-model support. No account required.",
};

export default function AiDemoPage() {
  return <AiDemo />;
}
