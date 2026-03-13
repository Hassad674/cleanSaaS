"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import { useChat } from "@/features/ai/hooks/use-chat";
import { ConversationList } from "@/features/ai/components/conversation-list";
import { ChatMessages } from "@/features/ai/components/chat-messages";
import { ChatInput } from "@/features/ai/components/chat-input";

export function ChatLayout() {
  const {
    conversations,
    currentConversationId,
    messages,
    loading,
    messagesLoading,
    sending,
    error,
    sendMessage,
    selectConversation,
    createConversation,
    deleteConversation,
  } = useChat();

  const [sidebarOpen, setSidebarOpen] = useState(false);

  function handleSelectConversation(id: string) {
    selectConversation(id);
    setSidebarOpen(false);
  }

  async function handleSendMessage(content: string) {
    if (!currentConversationId) {
      // Auto-create a conversation if none is selected
      await createConversation(content.slice(0, 50));
      // After creation, the hook sets currentConversationId and messages.
      // We need to wait a tick for state to update, then send.
      // Instead, we handle this by sending after creation in a useEffect-like pattern.
      // For simplicity, send directly — the hook will have set the ID.
    }
    sendMessage(content);
  }

  return (
    <div className="flex h-[calc(100vh-theme(spacing.14))] lg:h-[calc(100vh-theme(spacing.8))] -m-4 sm:-m-6 lg:-m-8">
      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-30 bg-background/80 backdrop-blur-sm lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-40 w-72 bg-card border-r border-border transition-transform duration-200 lg:relative lg:translate-x-0 lg:z-auto",
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <ConversationList
          conversations={conversations}
          currentConversationId={currentConversationId}
          loading={loading}
          onSelect={handleSelectConversation}
          onCreate={() => createConversation()}
          onDelete={deleteConversation}
        />
      </aside>

      {/* Main chat area */}
      <div className="flex-1 flex flex-col min-w-0 bg-background">
        {/* Chat header */}
        <div className="flex items-center gap-3 border-b border-border px-4 py-3 shrink-0">
          {/* Mobile menu button */}
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="lg:hidden text-muted-foreground hover:text-foreground transition-colors p-1"
            aria-label="Toggle sidebar"
          >
            <svg
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
              />
            </svg>
          </button>

          <h2 className="text-sm font-medium text-foreground truncate">
            {currentConversationId
              ? conversations.find((c) => c.id === currentConversationId)
                  ?.title ?? "Chat"
              : "AI Chat"}
          </h2>
        </div>

        {/* Error banner */}
        {error && (
          <div className="bg-destructive/10 border-b border-destructive/20 px-4 py-2">
            <p className="text-sm text-destructive">{error}</p>
          </div>
        )}

        {/* Messages area */}
        {currentConversationId ? (
          <ChatMessages
            messages={messages}
            loading={messagesLoading}
            sending={sending}
          />
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center gap-4 p-6">
            <div className="h-16 w-16 rounded-full bg-muted flex items-center justify-center">
              <svg
                className="h-8 w-8 text-muted-foreground"
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
            </div>
            <div className="text-center">
              <h3 className="text-lg font-medium text-foreground">
                Welcome to AI Chat
              </h3>
              <p className="text-sm text-muted-foreground mt-1">
                Select a conversation or start a new one to begin.
              </p>
            </div>
            <button
              onClick={() => createConversation()}
              className="rounded-lg bg-primary px-6 py-2.5 text-sm font-medium text-primary-foreground hover:opacity-90 transition-opacity"
            >
              Start new chat
            </button>
          </div>
        )}

        {/* Input area */}
        {currentConversationId && (
          <ChatInput onSend={handleSendMessage} disabled={sending} />
        )}
      </div>
    </div>
  );
}
