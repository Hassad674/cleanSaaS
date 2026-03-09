export type Notification = {
  id: string;
  title: string;
  body: string;
  channel: "email" | "in_app";
  read: boolean;
  created_at: string;
};
