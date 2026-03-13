"use client";

import { useState, useCallback } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type Role = "owner" | "admin" | "member";
type InviteStatus = "pending" | "accepted" | "declined";

interface MockTeam {
  id: string;
  name: string;
  slug: string;
  plan: string;
  memberCount: number;
  userRole: Role;
}

interface MockMember {
  id: string;
  name: string;
  email: string;
  role: Role;
  inviteStatus: InviteStatus;
  joinedAt?: Date;
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

const MOCK_TEAMS: MockTeam[] = [
  {
    id: "t-1",
    name: "Acme Corp",
    slug: "acme-corp",
    plan: "Pro",
    memberCount: 5,
    userRole: "owner",
  },
  {
    id: "t-2",
    name: "Side Project",
    slug: "side-project",
    plan: "Free",
    memberCount: 2,
    userRole: "admin",
  },
];

const INITIAL_MEMBERS: Record<string, MockMember[]> = {
  "t-1": [
    {
      id: "m-1",
      name: "You",
      email: "you@acme.com",
      role: "owner",
      inviteStatus: "accepted",
      joinedAt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000),
    },
    {
      id: "m-2",
      name: "Alice Johnson",
      email: "alice@acme.com",
      role: "admin",
      inviteStatus: "accepted",
      joinedAt: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000),
    },
    {
      id: "m-3",
      name: "Bob Smith",
      email: "bob@acme.com",
      role: "admin",
      inviteStatus: "accepted",
      joinedAt: new Date(Date.now() - 45 * 24 * 60 * 60 * 1000),
    },
    {
      id: "m-4",
      name: "Carol Williams",
      email: "carol@acme.com",
      role: "member",
      inviteStatus: "accepted",
      joinedAt: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000),
    },
    {
      id: "m-5",
      name: "David Brown",
      email: "david@acme.com",
      role: "member",
      inviteStatus: "pending",
    },
  ],
  "t-2": [
    {
      id: "m-6",
      name: "You",
      email: "you@side.io",
      role: "admin",
      inviteStatus: "accepted",
      joinedAt: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000),
    },
    {
      id: "m-7",
      name: "Eva Martinez",
      email: "eva@side.io",
      role: "owner",
      inviteStatus: "accepted",
      joinedAt: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000),
    },
  ],
};

// ---------------------------------------------------------------------------
// Icons
// ---------------------------------------------------------------------------

function ArrowLeftIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
    </svg>
  );
}

function PlusIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
    </svg>
  );
}

function UsersIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
    </svg>
  );
}

function UserPlusIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 0110.374 21c-2.331 0-4.512-.645-6.374-1.766z" />
    </svg>
  );
}

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

function RoleBadge({ role }: { role: Role }) {
  const styles: Record<Role, string> = {
    owner: "bg-primary/10 text-primary",
    admin: "bg-accent text-accent-foreground",
    member: "bg-muted text-muted-foreground",
  };

  return (
    <span className={cn("inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium", styles[role])}>
      {role.charAt(0).toUpperCase() + role.slice(1)}
    </span>
  );
}

function InviteStatusBadge({ status }: { status: InviteStatus }) {
  const styles: Record<InviteStatus, string> = {
    pending: "bg-warning/10 text-warning",
    accepted: "bg-success/10 text-success",
    declined: "bg-destructive/10 text-destructive",
  };

  return (
    <span className={cn("inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium", styles[status])}>
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  );
}

function TeamAvatar({ name }: { name: string }) {
  const initials = name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  return (
    <div className="h-12 w-12 rounded-lg bg-primary/10 text-primary flex items-center justify-center font-semibold text-sm flex-shrink-0">
      {initials}
    </div>
  );
}

function MemberAvatar({ name }: { name: string }) {
  const initials = name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  return (
    <div className="h-8 w-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold flex-shrink-0">
      {initials}
    </div>
  );
}

function formatDate(date: Date): string {
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function Toast({ message, onDismiss }: { message: string; onDismiss: () => void }) {
  return (
    <div className="fixed bottom-6 right-6 z-50 bg-card border border-border rounded-xl shadow-lg px-4 py-3 flex items-center gap-3 animate-in slide-in-from-bottom-4">
      <div className="h-6 w-6 rounded-full bg-success/10 text-success flex items-center justify-center">
        <CheckIcon className="h-3.5 w-3.5" />
      </div>
      <p className="text-sm text-foreground font-medium">{message}</p>
      <button
        type="button"
        onClick={onDismiss}
        className="text-muted-foreground hover:text-foreground transition-colors ml-2"
      >
        <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

export function TeamsDemo() {
  const [selectedTeamId, setSelectedTeamId] = useState<string | null>(null);
  const [membersMap, setMembersMap] = useState<Record<string, MockMember[]>>(INITIAL_MEMBERS);
  const [teams, setTeams] = useState<MockTeam[]>(MOCK_TEAMS);
  const [inviteDialogOpen, setInviteDialogOpen] = useState(false);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [inviteEmail, setInviteEmail] = useState("");
  const [inviteRole, setInviteRole] = useState<"admin" | "member">("member");
  const [newTeamName, setNewTeamName] = useState("");
  const [toast, setToast] = useState<string | null>(null);
  const [confirmRemoveId, setConfirmRemoveId] = useState<string | null>(null);

  const selectedTeam = teams.find((t) => t.id === selectedTeamId) ?? null;
  const currentMembers = selectedTeamId ? (membersMap[selectedTeamId] ?? []) : [];

  const showToast = useCallback((msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(null), 3000);
  }, []);

  // --- Create team ---
  const handleCreateTeam = () => {
    const trimmed = newTeamName.trim();
    if (!trimmed) return;

    const newTeam: MockTeam = {
      id: `t-${Date.now()}`,
      name: trimmed,
      slug: trimmed.toLowerCase().replace(/\s+/g, "-"),
      plan: "Free",
      memberCount: 1,
      userRole: "owner",
    };

    setTeams((prev) => [...prev, newTeam]);
    setMembersMap((prev) => ({
      ...prev,
      [newTeam.id]: [
        {
          id: `m-${Date.now()}`,
          name: "You",
          email: "you@example.com",
          role: "owner",
          inviteStatus: "accepted",
          joinedAt: new Date(),
        },
      ],
    }));
    setNewTeamName("");
    setCreateDialogOpen(false);
    showToast(`Team "${trimmed}" created`);
  };

  // --- Invite member ---
  const handleInvite = () => {
    if (!selectedTeamId || !inviteEmail.trim()) return;

    const newMember: MockMember = {
      id: `m-${Date.now()}`,
      name: inviteEmail.split("@")[0],
      email: inviteEmail.trim(),
      role: inviteRole,
      inviteStatus: "pending",
    };

    setMembersMap((prev) => ({
      ...prev,
      [selectedTeamId]: [...(prev[selectedTeamId] ?? []), newMember],
    }));

    setTeams((prev) =>
      prev.map((t) =>
        t.id === selectedTeamId ? { ...t, memberCount: t.memberCount + 1 } : t
      )
    );

    setInviteEmail("");
    setInviteRole("member");
    setInviteDialogOpen(false);
    showToast(`Invitation sent to ${newMember.email}`);
  };

  // --- Remove member ---
  const handleRemoveMember = (memberId: string) => {
    if (!selectedTeamId) return;

    setMembersMap((prev) => ({
      ...prev,
      [selectedTeamId]: (prev[selectedTeamId] ?? []).filter((m) => m.id !== memberId),
    }));

    setTeams((prev) =>
      prev.map((t) =>
        t.id === selectedTeamId ? { ...t, memberCount: Math.max(0, t.memberCount - 1) } : t
      )
    );

    setConfirmRemoveId(null);
    showToast("Member removed");
  };

  // --- Toggle role ---
  const handleToggleRole = (memberId: string) => {
    if (!selectedTeamId) return;

    setMembersMap((prev) => ({
      ...prev,
      [selectedTeamId]: (prev[selectedTeamId] ?? []).map((m) => {
        if (m.id !== memberId || m.role === "owner") return m;
        const nextRole = m.role === "admin" ? "member" : "admin";
        return { ...m, role: nextRole };
      }),
    }));

    showToast("Role updated");
  };

  const canManage = selectedTeam
    ? selectedTeam.userRole === "owner" || selectedTeam.userRole === "admin"
    : false;

  // =========================================================================
  // RENDER
  // =========================================================================

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 max-w-3xl">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeftIcon className="h-4 w-4" />
        Back to demos
      </Link>

      {/* TEAM LIST VIEW */}
      {!selectedTeamId && (
        <div className="space-y-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-foreground">Teams</h1>
              <p className="text-muted-foreground mt-1">
                Manage your teams and organizations.
              </p>
            </div>
            <button
              type="button"
              onClick={() => setCreateDialogOpen(true)}
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground text-sm font-medium rounded-lg hover:opacity-90 transition-opacity"
            >
              <PlusIcon className="h-4 w-4" />
              New team
            </button>
          </div>

          {/* Team cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {teams.map((team) => (
              <button
                key={team.id}
                type="button"
                onClick={() => setSelectedTeamId(team.id)}
                className="group bg-card border border-border rounded-xl p-5 shadow-sm hover:shadow-md hover:border-primary/50 transition-all duration-200 text-left flex flex-col"
              >
                <div className="flex items-start gap-4">
                  <TeamAvatar name={team.name} />
                  <div className="min-w-0 flex-1">
                    <h3 className="text-base font-semibold text-foreground truncate group-hover:text-primary transition-colors">
                      {team.name}
                    </h3>
                    <p className="text-sm text-muted-foreground mt-0.5">
                      {team.slug}
                    </p>
                  </div>
                </div>

                <div className="flex items-center justify-between mt-4 pt-3 border-t border-border">
                  <div className="flex items-center gap-2">
                    <RoleBadge role={team.userRole} />
                    <span className="text-xs text-muted-foreground capitalize">
                      {team.plan} plan
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                    <UsersIcon className="h-3.5 w-3.5" />
                    {team.memberCount} member{team.memberCount !== 1 ? "s" : ""}
                  </div>
                </div>
              </button>
            ))}
          </div>
        </div>
      )}

      {/* TEAM DETAIL VIEW */}
      {selectedTeamId && selectedTeam && (
        <div className="space-y-6">
          {/* Back to teams */}
          <button
            type="button"
            onClick={() => {
              setSelectedTeamId(null);
              setConfirmRemoveId(null);
            }}
            className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeftIcon className="h-4 w-4" />
            Back to teams
          </button>

          {/* Team info card */}
          <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
            <h2 className="text-base font-semibold text-foreground mb-4">
              Team settings
            </h2>
            <div className="flex items-center gap-4">
              <TeamAvatar name={selectedTeam.name} />
              <div>
                <h3 className="text-lg font-semibold text-foreground">
                  {selectedTeam.name}
                </h3>
                <div className="flex items-center gap-3 mt-1 text-xs text-muted-foreground">
                  <span>Slug: {selectedTeam.slug}</span>
                  <span>Plan: {selectedTeam.plan}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Members section */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">
                Members ({currentMembers.length})
              </h2>
              {canManage && (
                <button
                  type="button"
                  onClick={() => setInviteDialogOpen(true)}
                  className="inline-flex items-center gap-2 px-3 py-1.5 bg-primary text-primary-foreground text-sm font-medium rounded-lg hover:opacity-90 transition-opacity"
                >
                  <UserPlusIcon className="h-4 w-4" />
                  Invite
                </button>
              )}
            </div>

            {/* Members table */}
            <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden">
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
                      {canManage && (
                        <th className="text-right px-5 py-3 font-medium text-muted-foreground">
                          Actions
                        </th>
                      )}
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {currentMembers.map((member) => (
                      <tr
                        key={member.id}
                        className="hover:bg-muted/20 transition-colors"
                      >
                        <td className="px-5 py-3">
                          <div className="flex items-center gap-2.5">
                            <MemberAvatar name={member.name} />
                            <div>
                              <p className="text-foreground font-medium">
                                {member.name}
                              </p>
                              <p className="text-xs text-muted-foreground">
                                {member.email}
                              </p>
                            </div>
                          </div>
                        </td>
                        <td className="px-5 py-3">
                          <RoleBadge role={member.role} />
                        </td>
                        <td className="px-5 py-3">
                          <InviteStatusBadge status={member.inviteStatus} />
                        </td>
                        <td className="px-5 py-3 text-muted-foreground">
                          {member.joinedAt
                            ? formatDate(member.joinedAt)
                            : "---"}
                        </td>
                        {canManage && (
                          <td className="px-5 py-3">
                            {member.role !== "owner" && (
                              <div className="flex items-center justify-end gap-2">
                                <button
                                  type="button"
                                  onClick={() => handleToggleRole(member.id)}
                                  className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                                >
                                  {member.role === "admin"
                                    ? "Demote"
                                    : "Promote"}
                                </button>

                                {confirmRemoveId === member.id ? (
                                  <div className="flex items-center gap-1">
                                    <button
                                      type="button"
                                      onClick={() =>
                                        handleRemoveMember(member.id)
                                      }
                                      className="text-xs text-destructive font-medium hover:underline"
                                    >
                                      Confirm
                                    </button>
                                    <button
                                      type="button"
                                      onClick={() =>
                                        setConfirmRemoveId(null)
                                      }
                                      className="text-xs text-muted-foreground hover:text-foreground"
                                    >
                                      Cancel
                                    </button>
                                  </div>
                                ) : (
                                  <button
                                    type="button"
                                    onClick={() =>
                                      setConfirmRemoveId(member.id)
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
            </div>
          </div>
        </div>
      )}

      {/* CREATE TEAM DIALOG */}
      {createDialogOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div
            className="absolute inset-0 bg-background/80 backdrop-blur-sm"
            onClick={() => setCreateDialogOpen(false)}
          />
          <div className="relative bg-card border border-border rounded-xl shadow-lg w-full max-w-md mx-4 p-6">
            <h2 className="text-lg font-semibold text-foreground mb-1">
              Create a new team
            </h2>
            <p className="text-sm text-muted-foreground mb-6">
              Give your team a name to get started. You can change it later.
            </p>

            <form
              onSubmit={(e) => {
                e.preventDefault();
                handleCreateTeam();
              }}
            >
              <label
                htmlFor="demo-team-name"
                className="block text-sm font-medium text-foreground mb-1.5"
              >
                Team name
              </label>
              <input
                id="demo-team-name"
                type="text"
                value={newTeamName}
                onChange={(e) => setNewTeamName(e.target.value)}
                placeholder="e.g. Acme Corp"
                className="w-full bg-muted/50 border border-border rounded-lg px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                autoFocus
              />

              <div className="flex items-center justify-end gap-3 mt-6">
                <button
                  type="button"
                  onClick={() => setCreateDialogOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors rounded-lg"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={!newTeamName.trim()}
                  className={cn(
                    "px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity",
                    !newTeamName.trim() && "opacity-50 cursor-not-allowed"
                  )}
                >
                  Create team
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* INVITE MEMBER DIALOG */}
      {inviteDialogOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div
            className="absolute inset-0 bg-background/80 backdrop-blur-sm"
            onClick={() => setInviteDialogOpen(false)}
          />
          <div className="relative bg-card border border-border rounded-xl shadow-lg w-full max-w-md mx-4 p-6">
            <h2 className="text-lg font-semibold text-foreground mb-1">
              Invite a team member
            </h2>
            <p className="text-sm text-muted-foreground mb-6">
              Send an invitation by email. They&apos;ll receive a link to join
              your team.
            </p>

            <form
              onSubmit={(e) => {
                e.preventDefault();
                handleInvite();
              }}
            >
              <label
                htmlFor="demo-invite-email"
                className="block text-sm font-medium text-foreground mb-1.5"
              >
                Email address
              </label>
              <input
                id="demo-invite-email"
                type="email"
                value={inviteEmail}
                onChange={(e) => setInviteEmail(e.target.value)}
                placeholder="colleague@example.com"
                className="w-full bg-muted/50 border border-border rounded-lg px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                autoFocus
              />

              {/* Role selector */}
              <fieldset className="mt-4">
                <legend className="block text-sm font-medium text-foreground mb-2">
                  Role
                </legend>
                <div className="flex gap-3">
                  {(["member", "admin"] as const).map((r) => (
                    <button
                      key={r}
                      type="button"
                      onClick={() => setInviteRole(r)}
                      className={cn(
                        "flex-1 px-3 py-2 text-sm font-medium rounded-lg border transition-colors",
                        inviteRole === r
                          ? "border-primary bg-primary/10 text-primary"
                          : "border-border bg-muted/30 text-muted-foreground hover:text-foreground"
                      )}
                    >
                      {r === "admin" ? "Admin" : "Member"}
                    </button>
                  ))}
                </div>
                <p className="text-xs text-muted-foreground mt-2">
                  {inviteRole === "admin"
                    ? "Admins can manage members and team settings."
                    : "Members can view and collaborate within the team."}
                </p>
              </fieldset>

              <div className="flex items-center justify-end gap-3 mt-6">
                <button
                  type="button"
                  onClick={() => setInviteDialogOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors rounded-lg"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={!inviteEmail.trim()}
                  className={cn(
                    "px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity",
                    !inviteEmail.trim() && "opacity-50 cursor-not-allowed"
                  )}
                >
                  Send invite
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Toast */}
      {toast && <Toast message={toast} onDismiss={() => setToast(null)} />}

      {/* Footer note */}
      <p className="text-center text-sm text-muted-foreground mt-8">
        This is a demo with simulated data. In production, teams are managed
        server-side with role-based access control and invitation emails.
      </p>
    </div>
  );
}
