export type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

export type Conversation = {
  id: string;
  title: string;
  model: "claude" | "gpt" | "gemini";
  created_at: string;
  updated_at: string;
};
