"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";
import { ToastContainer } from "./toast";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export type ToastType = "success" | "error" | "warning" | "info";

export interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration: number;
  createdAt: number;
}

export interface ToastOptions {
  message: string;
  type?: ToastType;
  duration?: number;
}

interface ToastMethods {
  /** Show a toast with custom options */
  toast: (options: ToastOptions) => string;
  /** Shorthand for success toast */
  success: (message: string, duration?: number) => string;
  /** Shorthand for error toast */
  error: (message: string, duration?: number) => string;
  /** Shorthand for warning toast */
  warning: (message: string, duration?: number) => string;
  /** Shorthand for info toast */
  info: (message: string, duration?: number) => string;
  /** Dismiss a specific toast by id */
  dismiss: (id: string) => void;
  /** Dismiss all toasts */
  dismissAll: () => void;
}

interface ToastContextValue {
  toasts: Toast[];
  toast: ToastMethods;
}

// ---------------------------------------------------------------------------
// Context
// ---------------------------------------------------------------------------

const ToastContext = createContext<ToastContextValue | undefined>(undefined);

const DEFAULT_DURATION = 4000;
const MAX_VISIBLE = 5;

let toastCounter = 0;

function generateId(): string {
  toastCounter += 1;
  return `toast-${Date.now()}-${toastCounter}`;
}

// ---------------------------------------------------------------------------
// Provider
// ---------------------------------------------------------------------------

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);
  const timersRef = useRef<Map<string, ReturnType<typeof setTimeout>>>(
    new Map(),
  );

  const dismiss = useCallback((id: string) => {
    const timer = timersRef.current.get(id);
    if (timer) {
      clearTimeout(timer);
      timersRef.current.delete(id);
    }
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const dismissAll = useCallback(() => {
    timersRef.current.forEach((timer) => clearTimeout(timer));
    timersRef.current.clear();
    setToasts([]);
  }, []);

  const addToast = useCallback(
    (options: ToastOptions): string => {
      const id = generateId();
      const duration = options.duration ?? DEFAULT_DURATION;
      const newToast: Toast = {
        id,
        message: options.message,
        type: options.type ?? "info",
        duration,
        createdAt: Date.now(),
      };

      setToasts((prev) => {
        const next = [...prev, newToast];
        // If we exceed MAX_VISIBLE, remove the oldest
        if (next.length > MAX_VISIBLE) {
          const removed = next.splice(0, next.length - MAX_VISIBLE);
          removed.forEach((t) => {
            const timer = timersRef.current.get(t.id);
            if (timer) {
              clearTimeout(timer);
              timersRef.current.delete(t.id);
            }
          });
        }
        return next;
      });

      // Auto-dismiss after duration
      const timer = setTimeout(() => {
        timersRef.current.delete(id);
        setToasts((prev) => prev.filter((t) => t.id !== id));
      }, duration);
      timersRef.current.set(id, timer);

      return id;
    },
    [],
  );

  const toastMethods: ToastMethods = useMemo(
    () => ({
      toast: (options: ToastOptions) => addToast(options),
      success: (message: string, duration?: number) =>
        addToast({ message, type: "success", duration }),
      error: (message: string, duration?: number) =>
        addToast({ message, type: "error", duration }),
      warning: (message: string, duration?: number) =>
        addToast({ message, type: "warning", duration }),
      info: (message: string, duration?: number) =>
        addToast({ message, type: "info", duration }),
      dismiss,
      dismissAll,
    }),
    [addToast, dismiss, dismissAll],
  );

  const value = useMemo(
    () => ({ toasts, toast: toastMethods }),
    [toasts, toastMethods],
  );

  return (
    <ToastContext.Provider value={value}>
      {children}
      <ToastContainer toasts={toasts} onDismiss={dismiss} />
    </ToastContext.Provider>
  );
}

// ---------------------------------------------------------------------------
// Hook
// ---------------------------------------------------------------------------

export function useToast(): ToastMethods {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return context.toast;
}
