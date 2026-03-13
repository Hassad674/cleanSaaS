"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

/* -------------------------------------------------------------------------- */
/*  Mock data                                                                  */
/* -------------------------------------------------------------------------- */

type Role = "admin" | "user";
type Status = "active" | "inactive";
type Plan = "Free" | "Pro" | "Enterprise";

interface MockUser {
  id: string;
  name: string;
  email: string;
  role: Role;
  plan: Plan;
  status: Status;
  joinedAt: string;
}

const INITIAL_USERS: MockUser[] = [
  { id: "1", name: "Alice Martin", email: "alice@example.com", role: "admin", plan: "Enterprise", status: "active", joinedAt: "2025-01-12" },
  { id: "2", name: "John Doe", email: "john@example.com", role: "user", plan: "Pro", status: "active", joinedAt: "2025-02-03" },
  { id: "3", name: "Sarah Chen", email: "sarah@company.co", role: "user", plan: "Pro", status: "active", joinedAt: "2025-03-10" },
  { id: "4", name: "Mike Taylor", email: "mike@test.io", role: "user", plan: "Free", status: "inactive", joinedAt: "2024-11-22" },
  { id: "5", name: "Emma Wilson", email: "emma@startup.dev", role: "user", plan: "Enterprise", status: "active", joinedAt: "2025-01-28" },
  { id: "6", name: "Carlos Rivera", email: "carlos@agency.com", role: "admin", plan: "Pro", status: "active", joinedAt: "2024-10-15" },
  { id: "7", name: "Priya Sharma", email: "priya@design.co", role: "user", plan: "Free", status: "active", joinedAt: "2025-02-19" },
  { id: "8", name: "Liam O'Brien", email: "liam@corp.net", role: "user", plan: "Free", status: "inactive", joinedAt: "2024-12-05" },
];

const MONTHLY_REVENUE = [
  { month: "Oct", value: 4200 },
  { month: "Nov", value: 4800 },
  { month: "Dec", value: 5100 },
  { month: "Jan", value: 5900 },
  { month: "Feb", value: 6400 },
  { month: "Mar", value: 7296 },
];

interface ActivityItem {
  message: string;
  time: string;
}

const ACTIVITY_FEED: ActivityItem[] = [
  { message: "User john@example.com upgraded to Pro", time: "2h ago" },
  { message: "New user sarah@company.co registered", time: "5h ago" },
  { message: "Blog post 'Getting Started' published", time: "1d ago" },
  { message: "User mike@test.io deleted account", time: "2d ago" },
  { message: "Stripe webhook: invoice.paid for $19.00", time: "3d ago" },
];

/* -------------------------------------------------------------------------- */
/*  Icons (inline SVG to keep the component self-contained)                    */
/* -------------------------------------------------------------------------- */

function UsersIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
    </svg>
  );
}

function CreditCardIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z" />
    </svg>
  );
}

function CurrencyIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  );
}

function ServerIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z" />
    </svg>
  );
}

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
    </svg>
  );
}

function ArrowLeftIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
    </svg>
  );
}

function ActivityIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  );
}

/* -------------------------------------------------------------------------- */
/*  Stat card                                                                  */
/* -------------------------------------------------------------------------- */

interface StatCardProps {
  icon: React.ReactNode;
  label: string;
  value: string;
  change: string;
}

function StatCard({ icon, label, value, change }: StatCardProps) {
  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
      <div className="flex items-center justify-between mb-4">
        <div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
          {icon}
        </div>
        <span className="inline-flex items-center gap-1 text-xs font-medium text-emerald-600 bg-emerald-500/10 px-2 py-1 rounded-full">
          <svg className="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" />
          </svg>
          {change}
        </span>
      </div>
      <p className="text-sm text-muted-foreground mb-1">{label}</p>
      <p className="text-2xl font-bold text-foreground">{value}</p>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Bar chart (CSS-only)                                                       */
/* -------------------------------------------------------------------------- */

function RevenueChart() {
  const maxValue = Math.max(...MONTHLY_REVENUE.map((m) => m.value));

  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
      <h2 className="text-lg font-semibold text-foreground mb-6">Monthly Revenue</h2>
      <div className="flex items-end gap-3 sm:gap-4 h-48">
        {MONTHLY_REVENUE.map((item) => {
          const heightPercent = (item.value / maxValue) * 100;
          return (
            <div key={item.month} className="flex-1 flex flex-col items-center gap-2">
              <span className="text-xs font-medium text-muted-foreground">
                ${(item.value / 1000).toFixed(1)}k
              </span>
              <div className="w-full relative flex-1 flex items-end">
                <div
                  className="w-full bg-primary rounded-t-md transition-all duration-500 hover:opacity-80"
                  style={{ height: `${heightPercent}%` }}
                />
              </div>
              <span className="text-xs text-muted-foreground font-medium">{item.month}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  User table                                                                 */
/* -------------------------------------------------------------------------- */

function getInitials(name: string): string {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}

function formatJoinedDate(dateStr: string): string {
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(new Date(dateStr));
}

const PLAN_STYLES: Record<Plan, string> = {
  Free: "bg-muted text-muted-foreground",
  Pro: "bg-primary/10 text-primary",
  Enterprise: "bg-accent text-foreground",
};

interface UserTableProps {
  users: MockUser[];
  onToggleRole: (id: string) => void;
  onToggleStatus: (id: string) => void;
  search: string;
  onSearchChange: (value: string) => void;
}

function UserTable({ users, onToggleRole, onToggleStatus, search, onSearchChange }: UserTableProps) {
  return (
    <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden">
      {/* Header */}
      <div className="p-6 border-b border-border flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <h2 className="text-lg font-semibold text-foreground">User Management</h2>
        <div className="relative">
          <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search users..."
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            className="w-full sm:w-64 pl-9 pr-4 py-2 text-sm bg-muted border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-colors"
          />
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-muted/50">
              <th className="text-left px-6 py-3 font-medium text-muted-foreground">User</th>
              <th className="text-left px-6 py-3 font-medium text-muted-foreground hidden md:table-cell">Role</th>
              <th className="text-left px-6 py-3 font-medium text-muted-foreground hidden sm:table-cell">Plan</th>
              <th className="text-left px-6 py-3 font-medium text-muted-foreground">Status</th>
              <th className="text-left px-6 py-3 font-medium text-muted-foreground hidden lg:table-cell">Joined</th>
            </tr>
          </thead>
          <tbody>
            {users.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-muted-foreground">
                  No users match your search.
                </td>
              </tr>
            ) : (
              users.map((user) => (
                <tr
                  key={user.id}
                  className="border-b border-border last:border-b-0 hover:bg-muted/30 transition-colors"
                >
                  {/* User info */}
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="h-9 w-9 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold shrink-0">
                        {getInitials(user.name)}
                      </div>
                      <div className="min-w-0">
                        <p className="font-medium text-foreground truncate">{user.name}</p>
                        <p className="text-xs text-muted-foreground truncate">{user.email}</p>
                      </div>
                    </div>
                  </td>

                  {/* Role badge (clickable) */}
                  <td className="px-6 py-4 hidden md:table-cell">
                    <button
                      onClick={() => onToggleRole(user.id)}
                      className={cn(
                        "inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-medium transition-colors cursor-pointer",
                        user.role === "admin"
                          ? "bg-primary text-primary-foreground hover:opacity-80"
                          : "bg-muted text-muted-foreground hover:bg-muted/80"
                      )}
                      title="Click to toggle role"
                    >
                      {user.role}
                    </button>
                  </td>

                  {/* Plan badge */}
                  <td className="px-6 py-4 hidden sm:table-cell">
                    <span
                      className={cn(
                        "inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-medium",
                        PLAN_STYLES[user.plan]
                      )}
                    >
                      {user.plan}
                    </span>
                  </td>

                  {/* Status toggle */}
                  <td className="px-6 py-4">
                    <button
                      onClick={() => onToggleStatus(user.id)}
                      className="flex items-center gap-2 group cursor-pointer"
                      title="Click to toggle status"
                    >
                      <span
                        className={cn(
                          "h-2 w-2 rounded-full transition-colors",
                          user.status === "active" ? "bg-emerald-500" : "bg-muted-foreground/40"
                        )}
                      />
                      <span
                        className={cn(
                          "text-xs font-medium transition-colors",
                          user.status === "active"
                            ? "text-emerald-600"
                            : "text-muted-foreground"
                        )}
                      >
                        {user.status === "active" ? "Active" : "Inactive"}
                      </span>
                    </button>
                  </td>

                  {/* Joined date */}
                  <td className="px-6 py-4 hidden lg:table-cell">
                    <span className="text-muted-foreground">{formatJoinedDate(user.joinedAt)}</span>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Footer */}
      <div className="px-6 py-3 border-t border-border bg-muted/30">
        <p className="text-xs text-muted-foreground">
          Showing {users.length} of {INITIAL_USERS.length} users
        </p>
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Activity feed                                                              */
/* -------------------------------------------------------------------------- */

function ActivityFeed() {
  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
      <div className="flex items-center gap-2 mb-6">
        <ActivityIcon className="h-5 w-5 text-primary" />
        <h2 className="text-lg font-semibold text-foreground">Recent Activity</h2>
      </div>
      <div className="space-y-4">
        {ACTIVITY_FEED.map((item, i) => (
          <div key={i} className="flex items-start gap-3">
            <div className="mt-1.5 h-2 w-2 rounded-full bg-primary shrink-0" />
            <div className="min-w-0 flex-1">
              <p className="text-sm text-foreground leading-relaxed">{item.message}</p>
              <p className="text-xs text-muted-foreground mt-0.5">{item.time}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Main component                                                             */
/* -------------------------------------------------------------------------- */

export function AdminDemo() {
  const [users, setUsers] = useState<MockUser[]>(INITIAL_USERS);
  const [search, setSearch] = useState("");

  const filteredUsers = useMemo(() => {
    if (!search.trim()) return users;
    const q = search.toLowerCase();
    return users.filter(
      (u) =>
        u.name.toLowerCase().includes(q) ||
        u.email.toLowerCase().includes(q) ||
        u.role.includes(q) ||
        u.plan.toLowerCase().includes(q)
    );
  }, [users, search]);

  function handleToggleRole(id: string) {
    setUsers((prev) =>
      prev.map((u) =>
        u.id === id ? { ...u, role: u.role === "admin" ? "user" : "admin" } : u
      )
    );
  }

  function handleToggleStatus(id: string) {
    setUsers((prev) =>
      prev.map((u) =>
        u.id === id
          ? { ...u, status: u.status === "active" ? "inactive" : "active" }
          : u
      )
    );
  }

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeftIcon className="h-4 w-4" />
        Back to demos
      </Link>

      {/* Page header */}
      <div className="mb-8">
        <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
          Admin Dashboard
        </h1>
        <p className="text-muted-foreground mt-2">
          Overview of your SaaS platform. All data below is simulated for demo purposes.
        </p>
      </div>

      {/* Stats cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 mb-8">
        <StatCard
          icon={<UsersIcon className="h-5 w-5" />}
          label="Total Users"
          value="1,247"
          change="+12%"
        />
        <StatCard
          icon={<CreditCardIcon className="h-5 w-5" />}
          label="Active Subscriptions"
          value="384"
          change="+8%"
        />
        <StatCard
          icon={<CurrencyIcon className="h-5 w-5" />}
          label="Revenue (MRR)"
          value="$7,296"
          change="+15%"
        />
        <StatCard
          icon={<ServerIcon className="h-5 w-5" />}
          label="Storage Used"
          value="45.2 GB"
          change="+5%"
        />
      </div>

      {/* Chart + Activity row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6 mb-8">
        <div className="lg:col-span-2">
          <RevenueChart />
        </div>
        <div>
          <ActivityFeed />
        </div>
      </div>

      {/* User management table */}
      <UserTable
        users={filteredUsers}
        onToggleRole={handleToggleRole}
        onToggleStatus={handleToggleStatus}
        search={search}
        onSearchChange={setSearch}
      />

      {/* Footer note */}
      <p className="text-center text-xs text-muted-foreground mt-8">
        This is a self-contained demo with mock data. In production, this connects to your Go backend API.
      </p>
    </div>
  );
}
