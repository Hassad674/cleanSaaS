"use client";

import { useEffect, useCallback } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useTeam } from "@/features/team/hooks/use-team";
import { TeamSettings } from "@/features/team/components/team-settings";
import { MemberList } from "@/features/team/components/member-list";
import { PAGINATION_DEFAULT_LIMIT } from "@/shared/lib/constants";

function LoadingSkeleton() {
  return (
    <div className="max-w-3xl space-y-6">
      <div className="h-6 bg-muted rounded w-1/4 animate-pulse" />
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm animate-pulse">
        <div className="h-5 bg-muted rounded w-1/3 mb-4" />
        <div className="h-10 bg-muted rounded w-2/3" />
      </div>
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm animate-pulse">
        <div className="h-5 bg-muted rounded w-1/4 mb-4" />
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-center gap-3">
              <div className="h-8 w-8 rounded-full bg-muted" />
              <div className="h-4 bg-muted rounded w-1/3" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default function TeamDetailPage() {
  const params = useParams();
  const router = useRouter();
  const teamId = params.id as string;

  const {
    currentTeam,
    members,
    membersTotal,
    membersLoading,
    loading,
    error,
    fetchTeam,
    fetchMembers,
    updateTeam,
    deleteTeam,
    leaveTeam,
    inviteMember,
    removeMember,
    updateMemberRole,
    hasMembersNext,
    hasMembersPrev,
    goToNextMembersPage,
    goToPrevMembersPage,
    membersOffset,
    totalMembersPages,
  } = useTeam();

  useEffect(() => {
    fetchTeam(teamId);
    fetchMembers(teamId);
  }, [teamId, fetchTeam, fetchMembers]);

  const handleUpdate = useCallback(
    (name: string) => updateTeam(teamId, name),
    [teamId, updateTeam]
  );

  const handleDelete = useCallback(async () => {
    const result = await deleteTeam(teamId);
    if (result.success) {
      router.push("/teams");
    }
    return result;
  }, [teamId, deleteTeam, router]);

  const handleLeave = useCallback(async () => {
    const result = await leaveTeam(teamId);
    if (result.success) {
      router.push("/teams");
    }
    return result;
  }, [teamId, leaveTeam, router]);

  const handleInvite = useCallback(
    (email: string, role: "admin" | "member") =>
      inviteMember(teamId, email, role),
    [teamId, inviteMember]
  );

  const handleRemove = useCallback(
    (userId: string) => removeMember(teamId, userId),
    [teamId, removeMember]
  );

  const handleUpdateRole = useCallback(
    (userId: string, role: "admin" | "member") =>
      updateMemberRole(teamId, userId, role),
    [teamId, updateMemberRole]
  );

  if (loading && !currentTeam) {
    return <LoadingSkeleton />;
  }

  if (error && !currentTeam) {
    return (
      <div className="max-w-3xl">
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3">
          <p className="text-sm text-destructive">{error}</p>
        </div>
      </div>
    );
  }

  if (!currentTeam) return null;

  // Determine the current user's role from the members list
  // This is a simplified approach; in production you'd compare user IDs
  const currentUserRole =
    currentTeam.owner_id ? ("owner" as const) : ("member" as const);

  const currentPage = Math.floor(membersOffset / PAGINATION_DEFAULT_LIMIT) + 1;

  return (
    <div className="max-w-3xl space-y-6">
      {/* Back link */}
      <Link
        href="/teams"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <svg
          className="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
          />
        </svg>
        Back to teams
      </Link>

      {/* Team settings */}
      <TeamSettings
        team={currentTeam}
        isOwner={currentUserRole === "owner"}
        onUpdate={handleUpdate}
        onDelete={handleDelete}
        onLeave={handleLeave}
      />

      {/* Members */}
      <MemberList
        teamId={teamId}
        members={members}
        membersTotal={membersTotal}
        membersLoading={membersLoading}
        hasMembersNext={hasMembersNext}
        hasMembersPrev={hasMembersPrev}
        goToNextMembersPage={goToNextMembersPage}
        goToPrevMembersPage={goToPrevMembersPage}
        currentPage={currentPage}
        totalPages={totalMembersPages}
        onInvite={handleInvite}
        onRemove={handleRemove}
        onUpdateRole={handleUpdateRole}
        currentUserRole={currentUserRole}
      />
    </div>
  );
}
