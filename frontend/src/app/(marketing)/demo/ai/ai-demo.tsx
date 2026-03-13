"use client";

import { useState, useRef, useEffect, useCallback } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
}

interface Conversation {
  id: string;
  title: string;
  messages: Message[];
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

function uid(): string {
  return Math.random().toString(36).slice(2, 10);
}

const AI_RESPONSES = [
  "Great question! In CleanSaaS, every feature is fully independent and removable. The AI module uses a port/adapter pattern so you can swap between OpenAI, Gemini, or Claude with a single config change. No code modifications needed.",
  "Here is a quick example. To switch from OpenAI to Claude, you only need to change the provider in your environment config:\n\n```\nAI_PROVIDER=claude\nCLAUDE_API_KEY=sk-...\n```\n\nThe backend adapter handles the rest automatically.",
  "Sure, I can help with that! The backend is written in Go using the Chi router with a hexagonal architecture. Each feature has its own domain, ports, and adapters — making it easy to test, extend, or remove without touching other modules.",
  "Absolutely. The database layer uses pure SQL with golang-migrate for migrations. No ORM, no query builder. This gives you full control over performance and lets you write optimized queries for each use case.",
  "That is a great approach! For the frontend, CleanSaaS uses a feature-based architecture with Next.js 15. Each feature lives in its own folder under `features/` with components, actions, hooks, and types — completely self-contained.",
  "Happy to explain! The billing module integrates with Stripe for payments. It supports multiple plans, checkout sessions, customer portals, and webhook handling. Like every other module, it can be removed entirely if your project does not need billing.",
];

function initialConversations(): Conversation[] {
  return [
    {
      id: "conv-getting-started",
      title: "Getting Started",
      messages: [
        {
          id: uid(),
          role: "user",
          content: "What can you help me with?",
        },
        {
          id: uid(),
          role: "assistant",
          content:
            "I can help with coding, writing, analysis, math, and creative tasks. I\u2019m built on a swappable AI provider architecture \u2014 you can use Gemini, OpenAI, or Claude as the backend. What would you like to try?",
        },
      ],
    },
    {
      id: "conv-code-review",
      title: "Code Review",
      messages: [
        {
          id: uid(),
          role: "user",
          content: "Can you review Go code for best practices?",
        },
        {
          id: uid(),
          role: "assistant",
          content:
            "Of course! I can review Go code for idiomatic patterns, error handling, concurrency safety, and performance. Just paste your code and I\u2019ll provide detailed feedback with suggestions.",
        },
        {
          id: uid(),
          role: "user",
          content: "What about frontend React code?",
        },
        {
          id: uid(),
          role: "assistant",
          content:
            "Absolutely. I can review React and Next.js code for component structure, hook usage, performance optimizations, accessibility, and TypeScript best practices. I\u2019ll flag any anti-patterns and suggest improvements.",
        },
      ],
    },
    {
      id: "conv-recipe-ideas",
      title: "Recipe Ideas",
      messages: [
        {
          id: uid(),
          role: "user",
          content: "Give me a quick pasta recipe.",
        },
        {
          id: uid(),
          role: "assistant",
          content:
            "Here is a 15-minute aglio e olio:\n\n1. Boil 400g spaghetti in salted water until al dente\n2. Slice 6 cloves of garlic thinly\n3. Heat 1/3 cup olive oil, add garlic and a pinch of red pepper flakes\n4. Cook until garlic is golden (not brown!)\n5. Toss drained pasta in the oil, add pasta water as needed\n6. Finish with fresh parsley and parmesan\n\nSimple, fast, and delicious!",
        },
      ],
    },
  ];
}

// ---------------------------------------------------------------------------
// Icons (inline SVG to avoid external deps)
// ---------------------------------------------------------------------------

function IconPlus({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 4.5v15m7.5-7.5h-15"
      />
    </svg>
  );
}

function IconSend({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
      />
    </svg>
  );
}

function IconChat({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M2.25 12.76c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.076-4.076a1.526 1.526 0 011.037-.443 48.282 48.282 0 005.68-.494c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z"
      />
    </svg>
  );
}

function IconArrowLeft({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
      />
    </svg>
  );
}

function IconMenu({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
      />
    </svg>
  );
}

function IconX({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}

function IconSparkles({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z"
      />
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Typing indicator
// ---------------------------------------------------------------------------

function TypingIndicator() {
  return (
    <div className="flex items-start gap-3 max-w-[80%]">
      <div className="flex-shrink-0 h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center">
        <IconSparkles className="h-4 w-4 text-primary" />
      </div>
      <div className="bg-card border border-border rounded-2xl rounded-tl-md px-4 py-3">
        <div className="flex items-center gap-1.5">
          <span className="h-2 w-2 rounded-full bg-muted-foreground/60 animate-bounce [animation-delay:0ms]" />
          <span className="h-2 w-2 rounded-full bg-muted-foreground/60 animate-bounce [animation-delay:150ms]" />
          <span className="h-2 w-2 rounded-full bg-muted-foreground/60 animate-bounce [animation-delay:300ms]" />
        </div>
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Message bubble
// ---------------------------------------------------------------------------

function MessageBubble({ message }: { message: Message }) {
  const isUser = message.role === "user";

  return (
    <div
      className={cn(
        "flex items-start gap-3 animate-[fadeSlideUp_0.3s_ease-out]",
        isUser ? "flex-row-reverse" : "flex-row",
        isUser ? "max-w-[80%] ml-auto" : "max-w-[80%]"
      )}
    >
      {/* Avatar */}
      <div
        className={cn(
          "flex-shrink-0 h-8 w-8 rounded-full flex items-center justify-center text-xs font-medium",
          isUser
            ? "bg-primary text-primary-foreground"
            : "bg-primary/10 text-primary"
        )}
      >
        {isUser ? (
          "U"
        ) : (
          <IconSparkles className="h-4 w-4" />
        )}
      </div>

      {/* Bubble */}
      <div
        className={cn(
          "rounded-2xl px-4 py-2.5 text-sm leading-relaxed whitespace-pre-wrap",
          isUser
            ? "bg-primary text-primary-foreground rounded-tr-md"
            : "bg-card border border-border text-foreground rounded-tl-md"
        )}
      >
        {message.content}
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

export function AiDemo() {
  const [conversations, setConversations] = useState<Conversation[]>(
    initialConversations
  );
  const [activeConvId, setActiveConvId] = useState("conv-getting-started");
  const [input, setInput] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [responseIndex, setResponseIndex] = useState(0);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const activeConversation = conversations.find((c) => c.id === activeConvId);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [activeConversation?.messages, isTyping]);

  // Focus input on conversation switch
  useEffect(() => {
    inputRef.current?.focus();
  }, [activeConvId]);

  const streamResponse = useCallback(
    (convId: string, fullText: string) => {
      const words = fullText.split(" ");
      let current = "";
      const msgId = uid();
      let wordIndex = 0;

      // Add empty assistant message
      setConversations((prev) =>
        prev.map((c) =>
          c.id === convId
            ? {
                ...c,
                messages: [
                  ...c.messages,
                  { id: msgId, role: "assistant" as const, content: "" },
                ],
              }
            : c
        )
      );

      const interval = setInterval(() => {
        if (wordIndex < words.length) {
          current += (wordIndex > 0 ? " " : "") + words[wordIndex];
          wordIndex++;
          setConversations((prev) =>
            prev.map((c) =>
              c.id === convId
                ? {
                    ...c,
                    messages: c.messages.map((m) =>
                      m.id === msgId ? { ...m, content: current } : m
                    ),
                  }
                : c
            )
          );
        } else {
          clearInterval(interval);
          setIsTyping(false);
        }
      }, 40);
    },
    []
  );

  const handleSend = useCallback(() => {
    const trimmed = input.trim();
    if (!trimmed || isTyping || !activeConversation) return;

    // Add user message
    const userMsg: Message = { id: uid(), role: "user", content: trimmed };
    setConversations((prev) =>
      prev.map((c) =>
        c.id === activeConvId
          ? { ...c, messages: [...c.messages, userMsg] }
          : c
      )
    );
    setInput("");
    setIsTyping(true);

    // Pick next AI response and cycle
    const response = AI_RESPONSES[responseIndex % AI_RESPONSES.length];
    setResponseIndex((i) => i + 1);

    // Simulate delay before streaming starts
    const delay = 1000 + Math.random() * 1000;
    setTimeout(() => {
      streamResponse(activeConvId, response);
    }, delay);
  }, [input, isTyping, activeConversation, activeConvId, responseIndex, streamResponse]);

  const handleNewConversation = useCallback(() => {
    const newConv: Conversation = {
      id: `conv-${uid()}`,
      title: "New Chat",
      messages: [],
    };
    setConversations((prev) => [newConv, ...prev]);
    setActiveConvId(newConv.id);
    setSidebarOpen(false);
  }, []);

  const switchConversation = useCallback((convId: string) => {
    setActiveConvId(convId);
    setSidebarOpen(false);
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)]">
      {/* Top bar */}
      <div className="flex items-center gap-3 border-b border-border bg-card px-4 py-3">
        {/* Mobile sidebar toggle */}
        <button
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="lg:hidden flex-shrink-0 h-8 w-8 flex items-center justify-center rounded-lg hover:bg-muted transition-colors"
          aria-label={sidebarOpen ? "Close sidebar" : "Open sidebar"}
        >
          {sidebarOpen ? (
            <IconX className="h-5 w-5 text-foreground" />
          ) : (
            <IconMenu className="h-5 w-5 text-foreground" />
          )}
        </button>

        <Link
          href="/demo"
          className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <IconArrowLeft className="h-4 w-4" />
          <span className="hidden sm:inline">Back to demos</span>
        </Link>

        <div className="flex-1" />

        <div className="flex items-center gap-2">
          <IconSparkles className="h-5 w-5 text-primary" />
          <h1 className="text-sm font-semibold text-foreground">
            AI Chat Demo
          </h1>
        </div>

        <div className="flex-1" />

        {/* Spacer to balance layout */}
        <div className="w-8 lg:hidden" />
      </div>

      <div className="flex flex-1 overflow-hidden relative">
        {/* Sidebar overlay for mobile */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 bg-background/80 backdrop-blur-sm z-30 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
        )}

        {/* Sidebar */}
        <aside
          className={cn(
            "flex flex-col bg-card border-r border-border w-72 flex-shrink-0 transition-transform duration-200 ease-in-out z-40",
            "fixed inset-y-0 left-0 top-[calc(4rem+49px)] lg:static lg:translate-x-0",
            sidebarOpen ? "translate-x-0" : "-translate-x-full"
          )}
        >
          {/* New chat button */}
          <div className="p-3 border-b border-border">
            <button
              onClick={handleNewConversation}
              className="flex items-center justify-center gap-2 w-full rounded-lg border border-border bg-background px-3 py-2.5 text-sm font-medium text-foreground hover:bg-muted transition-colors"
            >
              <IconPlus className="h-4 w-4" />
              New conversation
            </button>
          </div>

          {/* Conversation list */}
          <nav className="flex-1 overflow-y-auto p-2 space-y-1">
            {conversations.map((conv) => (
              <button
                key={conv.id}
                onClick={() => switchConversation(conv.id)}
                className={cn(
                  "flex items-center gap-3 w-full rounded-lg px-3 py-2.5 text-sm transition-colors text-left",
                  conv.id === activeConvId
                    ? "bg-accent text-accent-foreground font-medium"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                )}
              >
                <IconChat className="h-4 w-4 flex-shrink-0" />
                <span className="truncate">{conv.title}</span>
              </button>
            ))}
          </nav>

          {/* Sidebar footer */}
          <div className="p-3 border-t border-border">
            <p className="text-xs text-muted-foreground text-center">
              Demo mode &mdash; no data is saved
            </p>
          </div>
        </aside>

        {/* Main chat area */}
        <main className="flex-1 flex flex-col min-w-0 bg-background">
          {/* Messages */}
          <div className="flex-1 overflow-y-auto px-4 py-6 space-y-4">
            {activeConversation && activeConversation.messages.length === 0 && (
              <div className="flex flex-col items-center justify-center h-full text-center gap-4">
                <div className="h-16 w-16 rounded-2xl bg-primary/10 flex items-center justify-center">
                  <IconSparkles className="h-8 w-8 text-primary" />
                </div>
                <div>
                  <h2 className="text-lg font-semibold text-foreground">
                    Start a conversation
                  </h2>
                  <p className="text-sm text-muted-foreground mt-1 max-w-sm">
                    Type a message below to see the AI chat in action. Responses
                    are simulated for this demo.
                  </p>
                </div>
              </div>
            )}

            {activeConversation?.messages.map((msg) => (
              <MessageBubble key={msg.id} message={msg} />
            ))}

            {isTyping && <TypingIndicator />}

            <div ref={messagesEndRef} />
          </div>

          {/* Input area */}
          <div className="border-t border-border bg-card px-4 py-3">
            <div className="flex items-center gap-2 max-w-3xl mx-auto">
              <input
                ref={inputRef}
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder={
                  isTyping ? "AI is responding..." : "Type a message..."
                }
                disabled={isTyping}
                className={cn(
                  "flex-1 rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground",
                  "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                  "disabled:opacity-50 disabled:cursor-not-allowed",
                  "transition-colors"
                )}
              />
              <button
                onClick={handleSend}
                disabled={isTyping || !input.trim()}
                className={cn(
                  "flex-shrink-0 h-10 w-10 rounded-lg flex items-center justify-center transition-all",
                  "bg-primary text-primary-foreground hover:opacity-90",
                  "disabled:opacity-50 disabled:cursor-not-allowed",
                  "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                )}
                aria-label="Send message"
              >
                <IconSend className="h-4 w-4" />
              </button>
            </div>
            <p className="text-xs text-muted-foreground text-center mt-2">
              This is a demo with simulated responses. In production, responses
              come from your configured AI provider.
            </p>
          </div>
        </main>
      </div>
    </div>
  );
}
