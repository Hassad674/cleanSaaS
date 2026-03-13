"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import { formatDate } from "@/shared/lib/utils";
import { RoleBadge } from "@/features/team/components/role-badge";
import { InviteMemberDialog } from "@/features/team/components/invite-member-dialog";
import type { TeamMember } from "@/features/team/types";

function MemberAvatar({ member }: { member: TeamMember }) {
  const label =
    member.invited_email
      ? member.invited_email[0].toUpperCase()
      : member.user_id.slice(0, 2).toUpperCase();

  return (
    <div className="h-8 w-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold flex-shrink-0">
      {label}
    </div>
  );
}

function InviteStatusBadge({
  status,
}: {
  status: TeamMember["invite_status"];
}) {
  const styles = {
    pending: "bg-warning/10 text-warning",
    accepted: "bg-success/10 text-success",
    declined: "bg-destructive/10 text-destructive",
  };

  return (
    <span
      className={`inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium ${styles[status]}`}
    >
      {status}
    </span>
  );
}

export function MemberList({
  teamId,
  members,
  membersTotal,
  membersLoading,
  hasMembersNext,
  hasMembersPrev,
  goToNextMembersPage,
  goToPrevMembersPage,
  currentPage,
  totalPages,
  onInvite,
  onRemove,
  onUpdateRole,
  currentUserRole,
}: {
  teamId: string;
  members: TeamMember[];
  membersTotal: number;
  membersLoading: boolean;
  hasMembersNext: boolean;
  hasMembersPrev: boolean;
  goToNextMembersPage: () => void;
  goToPrevMembersPage: () => void;
  currentPage: number;
  totalPages: number;
  onInvite: (
    email: string,
    role: "admin" | "member"
  ) => Promise<{ success: boolean; error: string | null }>;
  onRemove: (userId: string) => Promise<{ success: boolean; error: string | null }>;
  onUpdateRole: (
    userId: string,
    role: "admin" | "member"
  ) => Promise<{ success: boolean; error: string | null }>;
  currentUserRole?: "owner" | "admin" | "member";
}) {
  const [inviteDialogOpen, setInviteDialogOpen] = useState(false);
  const [confirmRemoveId, setConfirmRemoveId] = useState<string | null>(null);

  const canManageMembers =
    currentUserRole === "owner" || currentUserRole === "admin";

  const handleRemove = async (userId: string) => {
    await onRemove(userId);
    setConfirmRemoveId(null);
  };

  const handleRoleChange = async (
    member: TeamMember,
    newRole: "admin" | "member"
  ) => {
    await onUpdateRole(member.user_id, newRole);
  };

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">
          Members ({membersTotal})
        </h2>
        {canManageMembers && (
          <button
            type="button"
            onClick={() => setInviteDialogOpen(true)}
            className="inline-flex items-center gap-2 px-3 py-1.5 bg-primary text-primary-foreground text-sm font-medium rounded-lg hover:opacity-90 transition-opacity"
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
                d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 0110.374 21c-2.331 0-4.512-.645-6.374-1.766z"
              />
            </svg>
            Invite
          </button>
        )}
      </div>

      {/* Members table */}
      <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden">
        {membersLoading ? (
          <div className="p-6 space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center gap-3 animate-pulse">
                <div className="h-8 w-8 rounded-full bg-muted" />
                <div className="flex-1 space-y-1">
                  <div className="h-4 bg-muted rounded w-1/3" />
                  <div className="h-3 bg-muted rounded w-1/4" />
                </div>
              </div>
            ))}
          </div>
        ) : members.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 px-4">
            <h3 className="text-base font-medium text-foreground mb-1">
              No members yet
            </h3>
            <p className="text-sm text-muted-foreground text-center">
              Invite team members to get started.
            </p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border bg-muted/30">
                    <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                      Member
                    </th>
                    <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                      Role
                    </th>
                    <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                      Status
                    </th>
                    <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                      Joined
                    </th>
                    {canManageMembers && (
                      <th className="text-right px-5 py-3 font-medium text-muted-foreground">
                        Actions
                      </th>
                    )}
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {members.map((member) => (
                    <tr
                      key={member.id}
                      className="hover:bg-muted/20 transition-colors"
                    >
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-2.5">
                          <MemberAvatar member={member} />
                          <div>
                            <p className="text-foreground font-medium">
                              {member.invited_email ?? member.user_id.slice(0, 8)}
                            </p>
                            {member.invited_email && (
                              <p className="text-xs text-muted-foreground">
                                {member.invited_email}
                              </p>
                            )}
                          </div>
                        </div>
                      </td>
                      <td className="px-5 py-3">
                        <RoleBadge role={member.role} />
                      </td>
                      <td className="px-5 py-3">
                        <InviteStatusBadge status={member.invite_status} />
                      </td>
                      <td className="px-5 py-3 text-muted-foreground">
                        {member.joined_at
                          ? formatDate(member.joined_at)
                          : "---"}
                      </td>
                      {canManageMembers && (
                        <td className="px-5 py-3">
                          {member.role !== "owner" && (
                            <div className="flex items-center justify-end gap-2">
                              {/* Role toggle */}
                              <button
                                type="button"
                                onClick={() =>
                                  handleRoleChange(
                                    member,
                                    member.role === "admin"
                                      ? "member"
                                      : "admin"
                                  )
                                }
                                className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                                title={
                                  member.role === "admin"
                                    ? "Demote to member"
                                    : "Promote to admin"
                                }
                              >
                                {member.role === "admin"
                                  ? "Demote"
                                  : "Promote"}
                              </button>

                              {/* Remove */}
                              {confirmRemoveId === member.user_id ? (
                                <div className="flex items-center gap-1">
                                  <button
                                    type="button"
                                    onClick={() =>
                                      handleRemove(member.user_id)
                                    }
                                    className="text-xs text-destructive font-medium hover:underline"
                                  >
                                    Confirm
                                  </button>
                                  <button
                                    type="button"
                                    onClick={() => setConfirmRemoveId(null)}
                                    className="text-xs text-muted-foreground hover:text-foreground"
                                  >
                                    Cancel
                                  </button>
                                </div>
                              ) : (
                                <button
                                  type="button"
                                  onClick={() =>
                                    setConfirmRemoveId(member.user_id)
                                  }
                                  className="text-xs text-destructive hover:underline"
                                >
                                  Remove
                                </button>
                              )}
                            </div>
                          )}
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {(hasMembersPrev || hasMembersNext) && (
              <div className="flex items-center justify-between px-5 py-3 border-t border-border bg-muted/30">
                <button
                  type="button"
                  onClick={goToPrevMembersPage}
                  disabled={!hasMembersPrev}
                  className={cn(
                    "text-sm font-medium transition-colors",
                    hasMembersPrev
                      ? "text-foreground hover:text-primary"
                      : "text-muted-foreground/40 cursor-not-allowed"
                  )}
                >
                  Previous
                </button>
                <span className="text-xs text-muted-foreground">
                  Page {currentPage} of {totalPages}
                </span>
                <button
                  type="button"
                  onClick={goToNextMembersPage}
                  disabled={!hasMembersNext}
                  className={cn(
                    "text-sm font-medium transition-colors",
                    hasMembersNext
                      ? "text-foreground hover:text-primary"
                      : "text-muted-foreground/40 cursor-not-allowed"
                  )}
                >
                  Next
                </button>
              </div>
            )}
          </>
        )}
      </div>

      {/* Invite dialog */}
      <InviteMemberDialog
        open={inviteDialogOpen}
        onClose={() => setInviteDialogOpen(false)}
        onSubmit={(email, role) => onInvite(email, role)}
      />
    </div>
  );
}
