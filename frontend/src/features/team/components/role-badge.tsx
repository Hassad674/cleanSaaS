"use client";

import { cn } from "@/shared/lib/utils";

type Role = "owner" | "admin" | "member";

const roleStyles: Record<Role, string> = {
  owner: "bg-primary/10 text-primary",
  admin: "bg-accent text-accent-foreground",
  member: "bg-muted text-muted-foreground",
};

const roleLabels: Record<Role, string> = {
  owner: "Owner",
  admin: "Admin",
  member: "Member",
};

export function RoleBadge({
  role,
  className,
}: {
  role: Role;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium",
        roleStyles[role],
        className
      )}
    >
      {roleLabels[role]}
    </span>
  );
}
