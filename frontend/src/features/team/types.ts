export interface Team {
  id: string;
  name: string;
  slug: string;
  owner_id: string;
  avatar_url: string;
  plan: string;
  max_members: number;
  created_at: string;
}

export interface TeamMember {
  id: string;
  team_id: string;
  user_id: string;
  role: "owner" | "admin" | "member";
  invited_email?: string;
  invite_status: "pending" | "accepted" | "declined";
  joined_at?: string;
  created_at: string;
}

export type TeamListResponse = {
  teams: Team[];
  total: number;
};

export type TeamMemberListResponse = {
  members: TeamMember[];
  total: number;
  offset: number;
  limit: number;
};

export type TeamResponse = {
  team: Team;
};

export type InviteResponse = {
  message: string;
};

export type AcceptInviteResponse = {
  message: string;
};
