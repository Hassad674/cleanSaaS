"use client";

import { useState } from "react";
import { useToast } from "@/shared/components/toast-provider";
import type { ToastType } from "@/shared/components/toast-provider";

const typeLabels: { type: ToastType; label: string }[] = [
  { type: "success", label: "Success" },
  { type: "error", label: "Error" },
  { type: "warning", label: "Warning" },
  { type: "info", label: "Info" },
];

const defaultMessages: Record<ToastType, string> = {
  success: "Changes saved successfully!",
  error: "Something went wrong. Please try again.",
  warning: "Your session will expire in 5 minutes.",
  info: "A new version is available.",
};

const typeButtonStyles: Record<ToastType, string> = {
  success:
    "bg-success text-success-foreground hover:opacity-90",
  error:
    "bg-destructive text-destructive-foreground hover:opacity-90",
  warning:
    "bg-warning text-warning-foreground hover:opacity-90",
  info:
    "bg-primary text-primary-foreground hover:opacity-90",
};

export function ToastDemo() {
  const toast = useToast();
  const [customMessage, setCustomMessage] = useState("");
  const [duration, setDuration] = useState(4000);
  const [selectedType, setSelectedType] = useState<ToastType>("info");

  function handleQuickToast(type: ToastType) {
    toast[type](defaultMessages[type]);
  }

  function handleCustomToast() {
    const message = customMessage.trim() || "Hello from the toast system!";
    toast.toast({ message, type: selectedType, duration });
  }

  function handleFloodTest() {
    toast.success("First toast");
    setTimeout(() => toast.error("Second toast"), 200);
    setTimeout(() => toast.warning("Third toast"), 400);
    setTimeout(() => toast.info("Fourth toast"), 600);
    setTimeout(() => toast.success("Fifth toast"), 800);
    setTimeout(() => toast.error("Sixth toast (oldest removed)"), 1000);
  }

  return (
    <div className="space-y-10">
      {/* Quick triggers */}
      <section>
        <h2 className="text-xl font-bold text-foreground mb-2">
          Quick triggers
        </h2>
        <p className="text-sm text-muted-foreground mb-4">
          Click any button to fire a toast with default message.
        </p>
        <div className="flex flex-wrap gap-3">
          {typeLabels.map(({ type, label }) => (
            <button
              key={type}
              type="button"
              onClick={() => handleQuickToast(type)}
              className={[
                "px-4 py-2 rounded-lg text-sm font-medium transition-opacity",
                typeButtonStyles[type],
              ].join(" ")}
            >
              {label}
            </button>
          ))}
        </div>
      </section>

      {/* Custom toast builder */}
      <section className="bg-card border border-border rounded-xl p-6">
        <h2 className="text-xl font-bold text-foreground mb-4">
          Custom toast builder
        </h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          {/* Message */}
          <div className="space-y-2 sm:col-span-2">
            <label
              htmlFor="toast-message"
              className="block text-sm font-medium text-foreground"
            >
              Message
            </label>
            <input
              id="toast-message"
              type="text"
              value={customMessage}
              onChange={(e) => setCustomMessage(e.target.value)}
              placeholder="Enter your toast message..."
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>

          {/* Type selector */}
          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">
              Type
            </label>
            <div className="flex flex-wrap gap-2">
              {typeLabels.map(({ type, label }) => (
                <button
                  key={type}
                  type="button"
                  onClick={() => setSelectedType(type)}
                  className={[
                    "px-3 py-1.5 rounded-lg text-sm font-medium border transition-colors",
                    selectedType === type
                      ? "bg-primary text-primary-foreground border-primary"
                      : "bg-muted text-muted-foreground border-border hover:text-foreground",
                  ].join(" ")}
                >
                  {label}
                </button>
              ))}
            </div>
          </div>

          {/* Duration slider */}
          <div className="space-y-2">
            <label
              htmlFor="toast-duration"
              className="block text-sm font-medium text-foreground"
            >
              Duration: {(duration / 1000).toFixed(1)}s
            </label>
            <input
              id="toast-duration"
              type="range"
              min={1000}
              max={10000}
              step={500}
              value={duration}
              onChange={(e) => setDuration(Number(e.target.value))}
              className="w-full accent-primary"
            />
            <div className="flex justify-between text-xs text-muted-foreground">
              <span>1s</span>
              <span>10s</span>
            </div>
          </div>
        </div>

        <button
          type="button"
          onClick={handleCustomToast}
          className="mt-6 px-5 py-2.5 rounded-lg bg-primary text-primary-foreground text-sm font-medium hover:opacity-90 transition-opacity"
        >
          Show custom toast
        </button>
      </section>

      {/* Stress test */}
      <section>
        <h2 className="text-xl font-bold text-foreground mb-2">
          Queue stress test
        </h2>
        <p className="text-sm text-muted-foreground mb-4">
          Fires 6 toasts rapidly. Only 5 can be visible at once — the oldest
          gets evicted.
        </p>
        <button
          type="button"
          onClick={handleFloodTest}
          className="px-4 py-2 rounded-lg border border-border bg-muted text-foreground text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          Fire 6 toasts
        </button>
      </section>

      {/* Usage example */}
      <section className="bg-card border border-border rounded-xl p-6">
        <h2 className="text-xl font-bold text-foreground mb-4">
          Usage in code
        </h2>
        <pre className="bg-muted rounded-lg p-4 overflow-x-auto text-sm text-foreground font-mono leading-relaxed">
          <code>{`import { useToast } from "@/shared/components/toast-provider";

function MyComponent() {
  const toast = useToast();

  function handleSave() {
    try {
      // ... save logic
      toast.success("Changes saved!");
    } catch {
      toast.error("Failed to save changes.");
    }
  }

  return <button onClick={handleSave}>Save</button>;
}`}</code>
        </pre>
      </section>
    </div>
  );
}
