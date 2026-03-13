"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type ConnectionStatus = "disconnected" | "connecting" | "connected";
type EventType = "notification" | "chat" | "system";
type EventDirection = "received" | "sent";

interface WsEvent {
  id: string;
  type: EventType;
  payload: string;
  timestamp: Date;
  direction: EventDirection;
}

// ---------------------------------------------------------------------------
// Mock event templates
// ---------------------------------------------------------------------------

const NOTIFICATION_EVENTS = [
  "New comment on your post",
  "File upload complete",
  "Payment received: $19.00",
  "User invited to workspace",
  "Export ready for download",
];

const CHAT_EVENTS = [
  "AI response ready for 'Code Review'",
  "New message in General",
  "Typing indicator: Alex is writing...",
  "AI response ready for 'Architecture Plan'",
  "New message in #engineering",
];

const SYSTEM_EVENTS = [
  "Server health: OK",
  "Connected users: 42",
  "Memory usage: 67%",
  "Deployment v2.4.1 complete",
  "Cache cleared successfully",
];

const EVENT_TEMPLATES: Record<EventType, string[]> = {
  notification: NOTIFICATION_EVENTS,
  chat: CHAT_EVENTS,
  system: SYSTEM_EVENTS,
};

const EVENT_TYPES: EventType[] = ["notification", "chat", "system"];

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function randomBetween(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function formatTimestamp(date: Date): string {
  return date.toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  if (m === 0) return `${s}s`;
  return `${m}m ${s.toString().padStart(2, "0")}s`;
}

function generateEvent(direction: EventDirection, type?: EventType, payload?: string): WsEvent {
  const eventType = type ?? EVENT_TYPES[randomBetween(0, 2)];
  const templates = EVENT_TEMPLATES[eventType];
  const eventPayload = payload ?? templates[randomBetween(0, templates.length - 1)];
  return {
    id: `evt-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    type: eventType,
    payload: eventPayload,
    timestamp: new Date(),
    direction,
  };
}

function eventToJson(event: WsEvent): string {
  return JSON.stringify(
    {
      type: event.type,
      payload: { message: event.payload },
    },
    null,
    2,
  );
}

// ---------------------------------------------------------------------------
// Icons
// ---------------------------------------------------------------------------

function ArrowLeftIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
    </svg>
  );
}

function SignalIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M9.348 14.652a3.75 3.75 0 010-5.304m5.304 0a3.75 3.75 0 010 5.304m-7.425 2.121a6.75 6.75 0 010-9.546m9.546 0a6.75 6.75 0 010 9.546M5.106 18.894c-3.808-3.807-3.808-9.98 0-13.788m13.788 0c3.808 3.807 3.808 9.98 0 13.788M12 12h.008v.008H12V12zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
    </svg>
  );
}

function SendIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" />
    </svg>
  );
}

function ChartIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
    </svg>
  );
}

function CubeIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" />
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

function StatusDot({ status }: { status: ConnectionStatus }) {
  return (
    <span
      className={cn(
        "inline-block h-2.5 w-2.5 rounded-full",
        status === "connected" && "bg-success",
        status === "connecting" && "bg-warning animate-pulse",
        status === "disconnected" && "bg-destructive",
      )}
    />
  );
}

function TypeBadge({ type }: { type: EventType }) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md px-1.5 py-0.5 text-[11px] font-semibold uppercase tracking-wide",
        type === "notification" && "bg-primary/10 text-primary",
        type === "chat" && "bg-success/10 text-success",
        type === "system" && "bg-warning/10 text-warning",
      )}
    >
      {type}
    </span>
  );
}

function DirectionBadge({ direction }: { direction: EventDirection }) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wide",
        direction === "received" && "bg-muted text-muted-foreground",
        direction === "sent" && "bg-accent text-accent-foreground",
      )}
    >
      {direction === "received" ? "IN" : "OUT"}
    </span>
  );
}

// ---------------------------------------------------------------------------
// Connection Status Panel
// ---------------------------------------------------------------------------

function ConnectionStatusPanel({
  status,
  uptimeSeconds,
  onToggle,
}: {
  status: ConnectionStatus;
  uptimeSeconds: number;
  onToggle: () => void;
}) {
  const statusLabel =
    status === "connected"
      ? "Connected"
      : status === "connecting"
        ? "Connecting..."
        : "Disconnected";

  return (
    <div className="bg-card border border-border rounded-xl p-4 sm:p-5 shadow-sm">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
            <SignalIcon className="h-5 w-5 text-foreground" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <StatusDot status={status} />
              <span className="text-sm font-semibold text-foreground">{statusLabel}</span>
            </div>
            {status === "connected" && (
              <p className="text-xs text-muted-foreground mt-0.5">
                Connected for {formatDuration(uptimeSeconds)}
              </p>
            )}
            {status === "disconnected" && (
              <p className="text-xs text-muted-foreground mt-0.5">
                ws://localhost:8081/ws
              </p>
            )}
          </div>
        </div>

        <button
          type="button"
          onClick={onToggle}
          disabled={status === "connecting"}
          className={cn(
            "rounded-lg px-4 py-2 text-sm font-medium transition-opacity disabled:opacity-50",
            status === "connected"
              ? "bg-destructive text-destructive-foreground hover:opacity-90"
              : "bg-primary text-primary-foreground hover:opacity-90",
          )}
        >
          {status === "connected" ? "Disconnect" : "Connect"}
        </button>
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Live Event Feed
// ---------------------------------------------------------------------------

function LiveEventFeed({
  events,
  newEventIds,
}: {
  events: WsEvent[];
  newEventIds: Set<string>;
}) {
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = 0;
    }
  }, [events.length]);

  return (
    <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden flex flex-col">
      <div className="flex items-center justify-between px-4 sm:px-5 py-3 border-b border-border">
        <h2 className="text-sm font-semibold text-foreground">Live Event Feed</h2>
        <span className="inline-flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1.5 text-[11px] font-semibold text-primary-foreground">
          {events.length}
        </span>
      </div>

      <div
        ref={scrollRef}
        className="flex-1 overflow-y-auto divide-y divide-border"
        style={{ maxHeight: "480px" }}
      >
        {events.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 px-4">
            <SignalIcon className="h-10 w-10 text-muted-foreground/50 mb-3" />
            <p className="text-sm text-muted-foreground">
              Connect to start receiving events
            </p>
          </div>
        ) : (
          events.map((event) => (
            <div
              key={event.id}
              className={cn(
                "px-4 sm:px-5 py-3 transition-all duration-300",
                newEventIds.has(event.id) && "animate-[slideIn_0.4s_ease-out_forwards]",
              )}
            >
              <div className="flex items-start gap-3">
                <span className="text-xs font-mono text-muted-foreground whitespace-nowrap mt-0.5">
                  {formatTimestamp(event.timestamp)}
                </span>
                <div className="flex items-center gap-1.5 flex-shrink-0 mt-0.5">
                  <TypeBadge type={event.type} />
                  <DirectionBadge direction={event.direction} />
                </div>
                <p className="text-sm text-foreground leading-snug min-w-0 break-words">
                  {event.payload}
                </p>
              </div>
              <div className="mt-1.5 ml-0 sm:ml-[4.5rem]">
                <pre className="text-[11px] font-mono text-muted-foreground bg-muted rounded-md px-2.5 py-1.5 overflow-x-auto">
                  {eventToJson(event)}
                </pre>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Send Test Message Panel
// ---------------------------------------------------------------------------

function SendMessagePanel({
  connected,
  onSend,
}: {
  connected: boolean;
  onSend: (type: EventType, message: string) => void;
}) {
  const [message, setMessage] = useState("");
  const [type, setType] = useState<EventType>("notification");
  const [lastSentJson, setLastSentJson] = useState<string | null>(null);

  const handleSend = () => {
    if (!message.trim() || !connected) return;
    onSend(type, message.trim());
    setLastSentJson(
      JSON.stringify({ type, payload: { message: message.trim() } }, null, 2),
    );
    setMessage("");
  };

  return (
    <div className="bg-card border border-border rounded-xl p-4 sm:p-5 shadow-sm">
      <div className="flex items-center gap-2 mb-3">
        <SendIcon className="h-4 w-4 text-foreground" />
        <h2 className="text-sm font-semibold text-foreground">Send Test Message</h2>
      </div>

      <div className="space-y-3">
        <div>
          <label htmlFor="msg-type" className="block text-xs font-medium text-muted-foreground mb-1">
            Type
          </label>
          <select
            id="msg-type"
            value={type}
            onChange={(e) => setType(e.target.value as EventType)}
            disabled={!connected}
            className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground disabled:opacity-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <option value="notification">notification</option>
            <option value="chat">chat</option>
            <option value="system">system</option>
          </select>
        </div>

        <div>
          <label htmlFor="msg-body" className="block text-xs font-medium text-muted-foreground mb-1">
            Message
          </label>
          <input
            id="msg-body"
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") handleSend();
            }}
            disabled={!connected}
            placeholder={connected ? "Type a message..." : "Connect first"}
            className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground disabled:opacity-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>

        <button
          type="button"
          onClick={handleSend}
          disabled={!connected || !message.trim()}
          className="w-full rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          Send
        </button>

        {lastSentJson && (
          <div>
            <p className="text-[11px] font-medium text-muted-foreground mb-1">Last sent:</p>
            <pre className="text-[11px] font-mono text-muted-foreground bg-muted rounded-md px-2.5 py-1.5 overflow-x-auto">
              {lastSentJson}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Connection Metrics Panel
// ---------------------------------------------------------------------------

function MetricsPanel({
  receivedCount,
  sentCount,
  latency,
  uptimeSeconds,
  connected,
  heartbeatPulse,
}: {
  receivedCount: number;
  sentCount: number;
  latency: number;
  uptimeSeconds: number;
  connected: boolean;
  heartbeatPulse: boolean;
}) {
  return (
    <div className="bg-card border border-border rounded-xl p-4 sm:p-5 shadow-sm">
      <div className="flex items-center gap-2 mb-4">
        <ChartIcon className="h-4 w-4 text-foreground" />
        <h2 className="text-sm font-semibold text-foreground">Connection Metrics</h2>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <MetricCard label="Received" value={receivedCount.toString()} />
        <MetricCard label="Sent" value={sentCount.toString()} />
        <MetricCard
          label="Latency"
          value={connected ? `${latency}ms` : "--"}
        />
        <MetricCard
          label="Uptime"
          value={connected ? formatDuration(uptimeSeconds) : "--"}
        />
      </div>

      {/* Heartbeat visualization */}
      <div className="mt-4 pt-3 border-t border-border">
        <div className="flex items-center gap-2">
          <span className="text-xs font-medium text-muted-foreground">Heartbeat</span>
          <div className="flex items-center gap-1">
            {connected ? (
              <>
                <span
                  className={cn(
                    "inline-block h-2 w-2 rounded-full bg-success transition-transform duration-300",
                    heartbeatPulse && "scale-150",
                  )}
                />
                <span className="text-[11px] font-mono text-muted-foreground">
                  ping/pong
                </span>
              </>
            ) : (
              <>
                <span className="inline-block h-2 w-2 rounded-full bg-muted" />
                <span className="text-[11px] font-mono text-muted-foreground">
                  inactive
                </span>
              </>
            )}
          </div>
        </div>

        {connected && (
          <div className="mt-2 flex items-center gap-0.5">
            {Array.from({ length: 20 }).map((_, i) => (
              <div
                key={i}
                className={cn(
                  "h-3 w-1 rounded-full transition-all duration-300",
                  heartbeatPulse && i >= 8 && i <= 12
                    ? "bg-success h-5"
                    : "bg-muted",
                )}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function MetricCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg bg-muted/50 px-3 py-2.5">
      <p className="text-[11px] font-medium text-muted-foreground">{label}</p>
      <p className="text-lg font-bold font-mono text-foreground mt-0.5">{value}</p>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Architecture Diagram
// ---------------------------------------------------------------------------

function ArchitectureDiagram() {
  return (
    <div className="bg-card border border-border rounded-xl p-4 sm:p-5 shadow-sm">
      <div className="flex items-center gap-2 mb-4">
        <CubeIcon className="h-4 w-4 text-foreground" />
        <h2 className="text-sm font-semibold text-foreground">Architecture</h2>
      </div>

      {/* Diagram */}
      <div className="flex items-center justify-center gap-0 overflow-x-auto py-4">
        <DiagramNode label="Browser" sub="WebSocket client" />
        <DiagramArrow />
        <DiagramNode label="WebSocket" sub="ws:// upgrade" highlighted />
        <DiagramArrow />
        <DiagramNode label="Go Hub" sub="pkg/ws/" highlighted />
        <DiagramArrow />
        <div className="flex flex-col gap-1.5">
          <DiagramNode label="Notifications" sub="feature" small />
          <DiagramNode label="AI Chat" sub="feature" small />
        </div>
      </div>

      {/* Explanation */}
      <div className="mt-4 pt-3 border-t border-border">
        <p className="text-xs text-muted-foreground leading-relaxed">
          The WebSocket hub is a reusable <code className="font-mono bg-muted px-1 py-0.5 rounded text-foreground text-[11px]">pkg/ws/</code> module.
          Any feature can broadcast real-time events by implementing the{" "}
          <code className="font-mono bg-muted px-1 py-0.5 rounded text-foreground text-[11px]">Broadcaster</code> interface.
          The hub manages connections, rooms, and fan-out — features only call{" "}
          <code className="font-mono bg-muted px-1 py-0.5 rounded text-foreground text-[11px]">hub.Broadcast(event)</code>.
        </p>
      </div>
    </div>
  );
}

function DiagramNode({
  label,
  sub,
  highlighted,
  small,
}: {
  label: string;
  sub: string;
  highlighted?: boolean;
  small?: boolean;
}) {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center rounded-lg border text-center flex-shrink-0",
        highlighted
          ? "border-primary bg-primary/5"
          : "border-border bg-muted/50",
        small ? "px-2.5 py-1.5 min-w-[80px]" : "px-3 py-2.5 min-w-[90px]",
      )}
    >
      <span className={cn("font-semibold text-foreground", small ? "text-[11px]" : "text-xs")}>
        {label}
      </span>
      <span className={cn("text-muted-foreground", small ? "text-[9px]" : "text-[10px]")}>
        {sub}
      </span>
    </div>
  );
}

function DiagramArrow() {
  return (
    <div className="flex items-center px-1 flex-shrink-0">
      <div className="w-4 sm:w-6 h-px bg-border" />
      <div className="w-0 h-0 border-t-[4px] border-t-transparent border-b-[4px] border-b-transparent border-l-[6px] border-l-border" />
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main Component
// ---------------------------------------------------------------------------

export function RealtimeDemo() {
  const [status, setStatus] = useState<ConnectionStatus>("disconnected");
  const [events, setEvents] = useState<WsEvent[]>([]);
  const [newEventIds, setNewEventIds] = useState<Set<string>>(new Set());
  const [uptimeSeconds, setUptimeSeconds] = useState(0);
  const [latency, setLatency] = useState(12);
  const [sentCount, setSentCount] = useState(0);
  const [heartbeatPulse, setHeartbeatPulse] = useState(false);

  const connectedRef = useRef(false);
  const eventIntervalRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const receivedCount = useMemo(
    () => events.filter((e) => e.direction === "received").length,
    [events],
  );

  // ---- Connection toggle ----
  const handleToggle = useCallback(() => {
    if (status === "connecting") return;

    if (status === "connected") {
      // Disconnect
      connectedRef.current = false;
      setStatus("disconnected");
      setUptimeSeconds(0);

      // Add a system disconnect event
      const disconnectEvent = generateEvent("received", "system", "Connection closed by client");
      setEvents((prev) => [disconnectEvent, ...prev]);
      setNewEventIds((prev) => new Set(prev).add(disconnectEvent.id));
      setTimeout(() => {
        setNewEventIds((prev) => {
          const next = new Set(prev);
          next.delete(disconnectEvent.id);
          return next;
        });
      }, 500);
    } else {
      // Start connecting
      setStatus("connecting");
      setTimeout(() => {
        connectedRef.current = true;
        setStatus("connected");
        setUptimeSeconds(0);

        // Add a system connect event
        const connectEvent = generateEvent("received", "system", "WebSocket connection established");
        setEvents((prev) => [connectEvent, ...prev]);
        setNewEventIds((prev) => new Set(prev).add(connectEvent.id));
        setTimeout(() => {
          setNewEventIds((prev) => {
            const next = new Set(prev);
            next.delete(connectEvent.id);
            return next;
          });
        }, 500);
      }, 1000);
    }
  }, [status]);

  // ---- Uptime timer ----
  useEffect(() => {
    if (status !== "connected") return;
    const interval = setInterval(() => {
      setUptimeSeconds((s) => s + 1);
    }, 1000);
    return () => clearInterval(interval);
  }, [status]);

  // ---- Latency simulation ----
  useEffect(() => {
    if (status !== "connected") return;
    const interval = setInterval(() => {
      setLatency(randomBetween(8, 24));
    }, 3000);
    return () => clearInterval(interval);
  }, [status]);

  // ---- Heartbeat pulse ----
  useEffect(() => {
    if (status !== "connected") return;
    const interval = setInterval(() => {
      setHeartbeatPulse(true);
      setTimeout(() => setHeartbeatPulse(false), 300);
    }, 2000);
    return () => clearInterval(interval);
  }, [status]);

  // ---- Auto-generate events when connected ----
  useEffect(() => {
    if (status !== "connected") {
      if (eventIntervalRef.current) {
        clearTimeout(eventIntervalRef.current);
        eventIntervalRef.current = null;
      }
      return;
    }

    function scheduleNext() {
      const delay = randomBetween(3000, 5000);
      eventIntervalRef.current = setTimeout(() => {
        if (!connectedRef.current) return;
        const event = generateEvent("received");
        setEvents((prev) => [event, ...prev]);
        setNewEventIds((prev) => new Set(prev).add(event.id));
        setTimeout(() => {
          setNewEventIds((prev) => {
            const next = new Set(prev);
            next.delete(event.id);
            return next;
          });
        }, 500);
        scheduleNext();
      }, delay);
    }

    scheduleNext();

    return () => {
      if (eventIntervalRef.current) {
        clearTimeout(eventIntervalRef.current);
        eventIntervalRef.current = null;
      }
    };
  }, [status]);

  // ---- Send handler ----
  const handleSend = useCallback(
    (type: EventType, message: string) => {
      if (status !== "connected") return;
      const event = generateEvent("sent", type, message);
      setEvents((prev) => [event, ...prev]);
      setSentCount((c) => c + 1);
      setNewEventIds((prev) => new Set(prev).add(event.id));
      setTimeout(() => {
        setNewEventIds((prev) => {
          const next = new Set(prev);
          next.delete(event.id);
          return next;
        });
      }, 500);
    },
    [status],
  );

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 max-w-6xl">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-6"
      >
        <ArrowLeftIcon className="h-4 w-4" />
        Back to demos
      </Link>

      {/* Page header */}
      <div className="mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold text-foreground">
          Real-time WebSocket
        </h1>
        <p className="text-sm text-muted-foreground mt-1">
          Simulated WebSocket connection with live events, metrics, and message sending.
          No backend required.
        </p>
      </div>

      {/* Connection status - full width */}
      <div className="mb-6">
        <ConnectionStatusPanel
          status={status}
          uptimeSeconds={uptimeSeconds}
          onToggle={handleToggle}
        />
      </div>

      {/* Main layout: feed + sidebar */}
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-6">
        {/* Event feed - main area */}
        <div className="min-w-0">
          <LiveEventFeed events={events} newEventIds={newEventIds} />
        </div>

        {/* Sidebar */}
        <div className="flex flex-col gap-6">
          <MetricsPanel
            receivedCount={receivedCount}
            sentCount={sentCount}
            latency={latency}
            uptimeSeconds={uptimeSeconds}
            connected={status === "connected"}
            heartbeatPulse={heartbeatPulse}
          />
          <SendMessagePanel
            connected={status === "connected"}
            onSend={handleSend}
          />
        </div>
      </div>

      {/* Architecture diagram - full width */}
      <div className="mt-6">
        <ArchitectureDiagram />
      </div>

      {/* Footer note */}
      <p className="text-center text-sm text-muted-foreground mt-6">
        This is a simulated demo. In production, events arrive via WebSocket at{" "}
        <code className="font-mono bg-muted px-1 py-0.5 rounded text-foreground text-[11px]">
          GET /ws?token=xxx
        </code>{" "}
        from the Go backend.
      </p>
    </div>
  );
}
