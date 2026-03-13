import type { Metadata } from "next";
import Link from "next/link";
import { ToastDemo } from "./toast-demo";

export const metadata: Metadata = {
  title: "Toast Notifications Demo",
  description:
    "Interactive demo of the global toast notification system — success, error, warning, and info toasts.",
};

export default function ToastDemoPage() {
  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 max-w-3xl">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <svg
          className="h-4 w-4"
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
        Back to demos
      </Link>

      {/* Header */}
      <div className="mb-12">
        <span className="inline-block bg-primary/10 text-primary text-sm font-medium px-3 py-1 rounded-full mb-4">
          Shared Component
        </span>
        <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
          Toast Notifications
        </h1>
        <p className="text-lg text-muted-foreground mt-3">
          A global toast system for user feedback. Use{" "}
          <code className="bg-muted px-1.5 py-0.5 rounded text-sm font-mono text-foreground">
            useToast()
          </code>{" "}
          from any client component to trigger success, error, warning, or info
          toasts.
        </p>
      </div>

      {/* Demo */}
      <ToastDemo />
    </div>
  );
}
