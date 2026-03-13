"use client";

import { useAuth } from "@/features/auth/hooks/use-auth";

export function DashboardHome() {
  const { user, logout } = useAuth({ required: true });

  if (!user) return null;

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-foreground mb-6">Dashboard</h1>
      <div className="bg-card border border-border rounded-xl p-6 space-y-4">
        <div>
          <p className="text-sm text-muted-foreground">Welcome back,</p>
          <p className="text-lg font-semibold text-foreground">{user.name}</p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-muted-foreground">Email</p>
            <p className="text-foreground">{user.email}</p>
          </div>
          <div>
            <p className="text-muted-foreground">Role</p>
            <p className="text-foreground capitalize">{user.role}</p>
          </div>
        </div>
        <div className="pt-2">
          <button
            onClick={logout}
            className="text-sm text-destructive hover:underline"
          >
            Log out
          </button>
        </div>
      </div>
    </div>
  );
}
