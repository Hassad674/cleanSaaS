"use client";

import { useEffect, useRef, useState } from "react";
import type { Toast, ToastType } from "./toast-provider";

// ---------------------------------------------------------------------------
// Icons (inline SVG to avoid external dependencies)
// ---------------------------------------------------------------------------

function SuccessIcon() {
  return (
    <svg
      className="h-5 w-5 shrink-0"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function ErrorIcon() {
  return (
    <svg
      className="h-5 w-5 shrink-0"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function WarningIcon() {
  return (
    <svg
      className="h-5 w-5 shrink-0"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
      />
    </svg>
  );
}

function InfoIcon() {
  return (
    <svg
      className="h-5 w-5 shrink-0"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z"
      />
    </svg>
  );
}

function CloseIcon() {
  return (
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
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Style mappings per toast type (using design tokens)
// ---------------------------------------------------------------------------

const typeStyles: Record<
  ToastType,
  { container: string; icon: string; progress: string }
> = {
  success: {
    container:
      "bg-success/10 border-success/30 text-foreground",
    icon: "text-success",
    progress: "bg-success",
  },
  error: {
    container:
      "bg-destructive/10 border-destructive/30 text-foreground",
    icon: "text-destructive",
    progress: "bg-destructive",
  },
  warning: {
    container:
      "bg-warning/10 border-warning/30 text-foreground",
    icon: "text-warning",
    progress: "bg-warning",
  },
  info: {
    container:
      "bg-primary/10 border-primary/30 text-foreground",
    icon: "text-primary",
    progress: "bg-primary",
  },
};

const typeIcons: Record<ToastType, () => React.JSX.Element> = {
  success: SuccessIcon,
  error: ErrorIcon,
  warning: WarningIcon,
  info: InfoIcon,
};

// ---------------------------------------------------------------------------
// Single Toast Item
// ---------------------------------------------------------------------------

interface ToastItemProps {
  toast: Toast;
  onDismiss: (id: string) => void;
}

function ToastItem({ toast, onDismiss }: ToastItemProps) {
  const [isExiting, setIsExiting] = useState(false);
  const [isEntering, setIsEntering] = useState(true);
  const progressRef = useRef<HTMLDivElement>(null);

  const styles = typeStyles[toast.type];
  const Icon = typeIcons[toast.type];

  // Entry animation
  useEffect(() => {
    const frame = requestAnimationFrame(() => {
      setIsEntering(false);
    });
    return () => cancelAnimationFrame(frame);
  }, []);

  // Progress bar animation
  useEffect(() => {
    const el = progressRef.current;
    if (!el) return;

    // Start at full width, then shrink to 0 over the toast duration
    el.style.transition = "none";
    el.style.width = "100%";

    // Force reflow so the browser registers the initial state
    el.getBoundingClientRect();

    el.style.transition = `width ${toast.duration}ms linear`;
    el.style.width = "0%";
  }, [toast.duration]);

  function handleDismiss() {
    setIsExiting(true);
    setTimeout(() => onDismiss(toast.id), 200);
  }

  return (
    <div
      role="alert"
      aria-live="assertive"
      className={[
        "relative overflow-hidden border rounded-lg shadow-lg",
        "w-80 max-w-[calc(100vw-2rem)]",
        "transition-all duration-200 ease-out",
        styles.container,
        isEntering
          ? "translate-x-full sm:translate-x-full opacity-0"
          : "translate-x-0 opacity-100",
        isExiting ? "opacity-0 translate-x-4 scale-95" : "",
      ].join(" ")}
    >
      {/* Content */}
      <div className="flex items-start gap-3 p-4 pr-10">
        <span className={styles.icon}>
          <Icon />
        </span>
        <p className="text-sm font-medium leading-snug break-words">
          {toast.message}
        </p>
      </div>

      {/* Close button */}
      <button
        type="button"
        onClick={handleDismiss}
        className="absolute top-3 right-3 rounded-md p-1 text-muted-foreground hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        aria-label="Dismiss notification"
      >
        <CloseIcon />
      </button>

      {/* Progress bar */}
      <div className="h-1 w-full bg-muted/30">
        <div
          ref={progressRef}
          className={["h-full rounded-full opacity-60", styles.progress].join(
            " ",
          )}
        />
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Toast Container (renders all visible toasts)
// ---------------------------------------------------------------------------

interface ToastContainerProps {
  toasts: Toast[];
  onDismiss: (id: string) => void;
}

export function ToastContainer({ toasts, onDismiss }: ToastContainerProps) {
  if (toasts.length === 0) return null;

  return (
    <div
      aria-label="Notifications"
      className={[
        "fixed z-[9999] flex flex-col gap-2 pointer-events-none",
        // Bottom-center on mobile, bottom-right on desktop
        "bottom-4 left-1/2 -translate-x-1/2 items-center",
        "sm:bottom-6 sm:right-6 sm:left-auto sm:translate-x-0 sm:items-end",
      ].join(" ")}
    >
      {toasts.map((t) => (
        <div key={t.id} className="pointer-events-auto">
          <ToastItem toast={t} onDismiss={onDismiss} />
        </div>
      ))}
    </div>
  );
}
