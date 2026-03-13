export type Notification = {
  id: string;
  type: string;
  title: string;
  message: string;
  read: boolean;
  created_at: string;
};

export type NotificationsResponse = {
  notifications: Notification[];
  total: number;
  page: number;
  limit: number;
};

export type UnreadCountResponse = {
  count: number;
};
