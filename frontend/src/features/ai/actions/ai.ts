"use server";

import { api } from "@/shared/lib/api";
import type { Conversation, Message } from "@/features/ai/types";

type ConversationsResponse = {
  conversations: Conversation[];
  total: number;
  page: number;
  limit: number;
};

type MessagesResponse = {
  messages: Message[];
};

type ConversationResponse = Conversation;

type MessageResponse = Message;

export async function getConversations(
  authToken: string,
  page: number = 1,
  limit: number = 20
) {
  return api<ConversationsResponse>(
    `/ai/conversations?page=${page}&limit=${limit}`,
    { token: authToken }
  );
}

export async function createConversation(
  authToken: string,
  title?: string
) {
  return api<ConversationResponse>("/ai/conversations", {
    method: "POST",
    body: { title: title ?? "New conversation" },
    token: authToken,
  });
}

export async function getMessages(
  conversationId: string,
  authToken: string
) {
  return api<MessagesResponse>(
    `/ai/conversations/${conversationId}/messages`,
    { token: authToken }
  );
}

export async function sendMessage(
  conversationId: string,
  content: string,
  authToken: string
) {
  return api<MessageResponse>(
    `/ai/conversations/${conversationId}/messages`,
    {
      method: "POST",
      body: { content },
      token: authToken,
    }
  );
}

export async function deleteConversation(
  conversationId: string,
  authToken: string
) {
  return api<{ message: string }>(
    `/ai/conversations/${conversationId}`,
    {
      method: "DELETE",
      token: authToken,
    }
  );
}
