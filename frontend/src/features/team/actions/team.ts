"use server";

import { api } from "@/shared/lib/api";
import type {
  TeamResponse,
  TeamListResponse,
  TeamMemberListResponse,
  InviteResponse,
  AcceptInviteResponse,
} from "@/features/team/types";

export async function createTeam(name: string, authToken: string) {
  return api<TeamResponse>("/teams", {
    method: "POST",
    token: authToken,
    body: { name },
  });
}

export async function getUserTeams(authToken: string) {
  return api<TeamListResponse>("/teams", {
    token: authToken,
  });
}

export async function getTeam(id: string, authToken: string) {
  return api<TeamResponse>(`/teams/${id}`, {
    token: authToken,
  });
}

export async function updateTeam(
  id: string,
  name: string,
  authToken: string
) {
  return api<TeamResponse>(`/teams/${id}`, {
    method: "PATCH",
    token: authToken,
    body: { name },
  });
}

export async function deleteTeam(id: string, authToken: string) {
  return api<{ message: string }>(`/teams/${id}`, {
    method: "DELETE",
    token: authToken,
  });
}

export async function inviteMember(
  teamId: string,
  email: string,
  role: "admin" | "member",
  authToken: string
) {
  return api<InviteResponse>(`/teams/${teamId}/members/invite`, {
    method: "POST",
    token: authToken,
    body: { email, role },
  });
}

export async function acceptInvite(
  inviteToken: string,
  authToken: string
) {
  return api<AcceptInviteResponse>("/teams/invites/accept", {
    method: "POST",
    token: authToken,
    body: { token: inviteToken },
  });
}

export async function getTeamMembers(
  teamId: string,
  offset: number = 0,
  limit: number = 20,
  authToken: string
) {
  const params = new URLSearchParams({
    offset: String(offset),
    limit: String(limit),
  });

  return api<TeamMemberListResponse>(
    `/teams/${teamId}/members?${params.toString()}`,
    {
      token: authToken,
    }
  );
}

export async function removeMember(
  teamId: string,
  userId: string,
  authToken: string
) {
  return api<{ message: string }>(`/teams/${teamId}/members/${userId}`, {
    method: "DELETE",
    token: authToken,
  });
}

export async function updateMemberRole(
  teamId: string,
  userId: string,
  role: "admin" | "member",
  authToken: string
) {
  return api<{ message: string }>(`/teams/${teamId}/members/${userId}/role`, {
    method: "PATCH",
    token: authToken,
    body: { role },
  });
}

export async function leaveTeam(teamId: string, authToken: string) {
  return api<{ message: string }>(`/teams/${teamId}/leave`, {
    method: "POST",
    token: authToken,
  });
}
