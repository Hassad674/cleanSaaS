"use client";

import Link from "next/link";
import { cn } from "@/shared/lib/utils";
import { RoleBadge } from "@/features/team/components/role-badge";
import type { Team } from "@/features/team/types";

function TeamAvatar({
  name,
  avatarUrl,
  className,
}: {
  name: string;
  avatarUrl?: string;
  className?: string;
}) {
  const initials = name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  if (avatarUrl) {
    return (
      <img
        src={avatarUrl}
        alt={name}
        className={cn("rounded-lg object-cover", className)}
      />
    );
  }

  return (
    <div
      className={cn(
        "rounded-lg bg-primary/10 text-primary flex items-center justify-center font-semibold",
        className
      )}
    >
      {initials}
    </div>
  );
}

export function TeamCard({
  team,
  memberCount,
  userRole,
}: {
  team: Team;
  memberCount?: number;
  userRole?: "owner" | "admin" | "member";
}) {
  return (
    <Link
      href={`/teams/${team.id}`}
      className="group bg-card border border-border rounded-xl p-5 shadow-sm hover:shadow-md hover:border-primary/50 transition-all duration-200 flex flex-col"
    >
      <div className="flex items-start gap-4">
        <TeamAvatar
          name={team.name}
          avatarUrl={team.avatar_url}
          className="h-12 w-12 text-sm flex-shrink-0"
        />
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
          {userRole && <RoleBadge role={userRole} />}
          <span className="text-xs text-muted-foreground capitalize">
            {team.plan} plan
          </span>
        </div>
        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          <svg
            className="h-3.5 w-3.5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={1.5}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"
            />
          </svg>
          {memberCount ?? team.max_members} member{(memberCount ?? team.max_members) !== 1 ? "s" : ""}
        </div>
      </div>
    </Link>
  );
}
