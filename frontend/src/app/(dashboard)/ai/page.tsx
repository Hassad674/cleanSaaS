import type { Metadata } from "next";

export const metadata: Metadata = { title: "AI Chat" };

export default function AIPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">AI Chat</h1>
      <p className="text-muted-foreground">AI conversation interface will be here.</p>
    </div>
  );
}
