"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { useAuth } from "@/shared/hooks/use-auth";
import {
  getConversations as getConversationsAction,
  createConversation as createConversationAction,
  getMessages as getMessagesAction,
  deleteConversation as deleteConversationAction,
} from "@/features/ai/actions/ai";
import type { Conversation, Message } from "@/features/ai/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

export function useChat() {
  const { getToken } = useAuth({ required: true });

  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [currentConversationId, setCurrentConversationId] = useState<
    string | null
  >(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);
  const [messagesLoading, setMessagesLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const abortControllerRef = useRef<AbortController | null>(null);

  const fetchConversations = useCallback(() => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);
    getConversationsAction(token, 1, 50).then((res) => {
      if (res.data) {
        setConversations(res.data.conversations ?? []);
      } else {
        setError(res.error ?? "Failed to load conversations");
      }
      setLoading(false);
    });
  }, [getToken]);

  useEffect(() => {
    fetchConversations();
  }, [fetchConversations]);

  const selectConversation = useCallback(
    (id: string) => {
      const token = getToken();
      if (!token) return;

      setCurrentConversationId(id);
      setMessagesLoading(true);
      setMessages([]);
      setError(null);

      getMessagesAction(id, token).then((res) => {
        if (res.data) {
          setMessages(res.data.messages ?? []);
        } else {
          setError(res.error ?? "Failed to load messages");
        }
        setMessagesLoading(false);
      });
    },
    [getToken]
  );

  const createConversation = useCallback(
    async (title?: string) => {
      const token = getToken();
      if (!token) return;

      setError(null);
      const res = await createConversationAction(token, title);

      if (res.data) {
        setConversations((prev) => [res.data!, ...prev]);
        setCurrentConversationId(res.data.id);
        setMessages([]);
      } else {
        setError(res.error ?? "Failed to create conversation");
      }
    },
    [getToken]
  );

  const sendMessage = useCallback(
    async (content: string) => {
      const token = getToken();
      if (!token || !currentConversationId || sending) return;

      // Cancel any in-flight stream
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      const userMessage: Message = { role: "user", content };
      setMessages((prev) => [...prev, userMessage]);
      setSending(true);
      setError(null);

      // Add a placeholder for the assistant response
      const assistantMessage: Message = { role: "assistant", content: "" };
      setMessages((prev) => [...prev, assistantMessage]);

      const controller = new AbortController();
      abortControllerRef.current = controller;

      try {
        const res = await fetch(
          `${API_URL}/ai/conversations/${currentConversationId}/stream`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ content }),
            signal: controller.signal,
          }
        );

        if (!res.ok) {
          const errorData = await res.json().catch(() => null);
          throw new Error(errorData?.error ?? "Failed to send message");
        }

        const reader = res.body?.getReader();
        if (!reader) {
          throw new Error("No response stream available");
        }

        const decoder = new TextDecoder();
        let buffer = "";

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });

          // Process SSE events from buffer
          const lines = buffer.split("\n");
          // Keep the last incomplete line in the buffer
          buffer = lines.pop() ?? "";

          for (const line of lines) {
            if (line.startsWith("event: done")) {
              // Stream is complete
              break;
            }
            if (line.startsWith("data: ")) {
              const chunk = line.slice(6);
              setMessages((prev) => {
                const updated = [...prev];
                const lastIdx = updated.length - 1;
                if (lastIdx >= 0 && updated[lastIdx].role === "assistant") {
                  updated[lastIdx] = {
                    ...updated[lastIdx],
                    content: updated[lastIdx].content + chunk,
                  };
                }
                return updated;
              });
            }
          }
        }

        // Update the conversation title in the list if it was a new conversation
        fetchConversations();
      } catch (err: unknown) {
        if (err instanceof DOMException && err.name === "AbortError") {
          // User cancelled, not an error
          return;
        }

        const errorMessage =
          err instanceof Error ? err.message : "Failed to send message";
        setError(errorMessage);

        // Remove the empty assistant placeholder on error
        setMessages((prev) => {
          const updated = [...prev];
          const lastIdx = updated.length - 1;
          if (
            lastIdx >= 0 &&
            updated[lastIdx].role === "assistant" &&
            updated[lastIdx].content === ""
          ) {
            updated.pop();
          }
          return updated;
        });
      } finally {
        setSending(false);
        abortControllerRef.current = null;
      }
    },
    [getToken, currentConversationId, sending, fetchConversations]
  );

  const deleteConversation = useCallback(
    async (id: string) => {
      const token = getToken();
      if (!token) return;

      setError(null);
      const res = await deleteConversationAction(id, token);

      if (res.error) {
        setError(res.error);
        return;
      }

      setConversations((prev) => prev.filter((c) => c.id !== id));

      if (currentConversationId === id) {
        setCurrentConversationId(null);
        setMessages([]);
      }
    },
    [getToken, currentConversationId]
  );

  const stopStreaming = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
      setSending(false);
    }
  }, []);

  return {
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
    stopStreaming,
  };
}
