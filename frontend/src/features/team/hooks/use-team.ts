"use client";

import { useState, useEffect, useCallback } from "react";
import { useAuth } from "@/shared/hooks/use-auth";
import {
  getUserTeams as getUserTeamsAction,
  getTeam as getTeamAction,
  createTeam as createTeamAction,
  updateTeam as updateTeamAction,
  deleteTeam as deleteTeamAction,
  getTeamMembers as getTeamMembersAction,
  inviteMember as inviteMemberAction,
  removeMember as removeMemberAction,
  updateMemberRole as updateMemberRoleAction,
  leaveTeam as leaveTeamAction,
} from "@/features/team/actions/team";
import { PAGINATION_DEFAULT_LIMIT } from "@/shared/lib/constants";
import type { Team, TeamMember } from "@/features/team/types";

export function useTeam() {
  const { getToken } = useAuth({ required: true });

  const [teams, setTeams] = useState<Team[]>([]);
  const [currentTeam, setCurrentTeam] = useState<Team | null>(null);
  const [members, setMembers] = useState<TeamMember[]>([]);
  const [membersTotal, setMembersTotal] = useState(0);
  const [membersOffset, setMembersOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [membersLoading, setMembersLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const limit = PAGINATION_DEFAULT_LIMIT;

  // ---- Teams list ----

  const fetchTeams = useCallback(() => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);
    getUserTeamsAction(token).then((res) => {
      if (res.data) {
        setTeams(res.data.teams ?? []);
      } else {
        setError(res.error ?? "Failed to load teams");
      }
      setLoading(false);
    });
  }, [getToken]);

  // ---- Single team ----

  const fetchTeam = useCallback(
    (id: string) => {
      const token = getToken();
      if (!token) return;

      setLoading(true);
      setError(null);
      getTeamAction(id, token).then((res) => {
        if (res.data) {
          setCurrentTeam(res.data.team);
        } else {
          setError(res.error ?? "Failed to load team");
        }
        setLoading(false);
      });
    },
    [getToken]
  );

  // ---- Members ----

  const fetchMembers = useCallback(
    (teamId: string) => {
      const token = getToken();
      if (!token) return;

      setMembersLoading(true);
      getTeamMembersAction(teamId, membersOffset, limit, token).then((res) => {
        if (res.data) {
          setMembers(res.data.members ?? []);
          setMembersTotal(res.data.total);
        } else {
          setError(res.error ?? "Failed to load members");
        }
        setMembersLoading(false);
      });
    },
    [getToken, membersOffset, limit]
  );

  // ---- CRUD operations ----

  const createTeam = useCallback(
    async (name: string) => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await createTeamAction(name, token);
      if (res.data) {
        fetchTeams();
        return { success: true, error: null };
      }
      return { success: false, error: res.error ?? "Failed to create team" };
    },
    [getToken, fetchTeams]
  );

  const updateTeam = useCallback(
    async (id: string, name: string) => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await updateTeamAction(id, name, token);
      if (res.data) {
        setCurrentTeam(res.data.team);
        fetchTeams();
        return { success: true, error: null };
      }
      return { success: false, error: res.error ?? "Failed to update team" };
    },
    [getToken, fetchTeams]
  );

  const deleteTeam = useCallback(
    async (id: string) => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await deleteTeamAction(id, token);
      if (res.data) {
        setCurrentTeam(null);
        fetchTeams();
        return { success: true, error: null };
      }
      return { success: false, error: res.error ?? "Failed to delete team" };
    },
    [getToken, fetchTeams]
  );

  const inviteMember = useCallback(
    async (teamId: string, email: string, role: "admin" | "member") => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await inviteMemberAction(teamId, email, role, token);
      if (res.data) {
        fetchMembers(teamId);
        return { success: true, error: null };
      }
      return { success: false, error: res.error ?? "Failed to invite member" };
    },
    [getToken, fetchMembers]
  );

  const removeMember = useCallback(
    async (teamId: string, userId: string) => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await removeMemberAction(teamId, userId, token);
      if (res.data) {
        fetchMembers(teamId);
        return { success: true, error: null };
      }
      return { success: false, error: res.error ?? "Failed to remove member" };
    },
    [getToken, fetchMembers]
  );

  const updateMemberRole = useCallback(
    async (teamId: string, userId: string, role: "admin" | "member") => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await updateMemberRoleAction(teamId, userId, role, token);
      if (res.data) {
        fetchMembers(teamId);
        return { success: true, error: null };
      }
      return {
        success: false,
        error: res.error ?? "Failed to update member role",
      };
    },
    [getToken, fetchMembers]
  );

  const leaveTeam = useCallback(
    async (teamId: string) => {
      const token = getToken();
      if (!token) return { success: false, error: "Not authenticated" };

      const res = await leaveTeamAction(teamId, token);
      if (res.data) {
        setCurrentTeam(null);
        fetchTeams();
        return { success: true, error: null };
      }
      return { success: false, error: res.error ?? "Failed to leave team" };
    },
    [getToken, fetchTeams]
  );

  // ---- Pagination ----

  const totalMembersPages = Math.ceil(membersTotal / limit);
  const hasMembersNext = membersOffset + limit < membersTotal;
  const hasMembersPrev = membersOffset > 0;

  const goToNextMembersPage = useCallback(() => {
    if (hasMembersNext) setMembersOffset((prev) => prev + limit);
  }, [hasMembersNext, limit]);

  const goToPrevMembersPage = useCallback(() => {
    if (hasMembersPrev) setMembersOffset((prev) => Math.max(0, prev - limit));
  }, [hasMembersPrev, limit]);

  // Fetch teams on mount
  useEffect(() => {
    fetchTeams();
  }, [fetchTeams]);

  return {
    // Teams
    teams,
    currentTeam,
    loading,
    error,
    fetchTeams,
    fetchTeam,
    createTeam,
    updateTeam,
    deleteTeam,
    leaveTeam,

    // Members
    members,
    membersTotal,
    membersLoading,
    fetchMembers,
    inviteMember,
    removeMember,
    updateMemberRole,

    // Pagination
    membersOffset,
    totalMembersPages,
    hasMembersNext,
    hasMembersPrev,
    goToNextMembersPage,
    goToPrevMembersPage,
  };
}
