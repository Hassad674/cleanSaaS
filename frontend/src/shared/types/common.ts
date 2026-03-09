export type User = {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  role: "admin" | "member";
  email_verified: boolean;
};

export type AuthResponse = {
  token: string;
  user: User;
};
