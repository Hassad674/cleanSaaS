"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type ReferralStatus = "pending" | "completed" | "rewarded";

interface MockReferral {
  id: string;
  referredName: string;
  status: ReferralStatus;
  createdAt: Date;
  completedAt?: Date;
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

const MOCK_CODE = "CLEAN2024";

const MOCK_REFERRALS: MockReferral[] = [
  {
    id: "r-1",
    referredName: "Alice Johnson",
    status: "rewarded",
    createdAt: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000),
    completedAt: new Date(Date.now() - 25 * 24 * 60 * 60 * 1000),
  },
  {
    id: "r-2",
    referredName: "Bob Smith",
    status: "completed",
    createdAt: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000),
    completedAt: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
  },
  {
    id: "r-3",
    referredName: "Carol Williams",
    status: "completed",
    createdAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
    completedAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
  },
  {
    id: "r-4",
    referredName: "David Brown",
    status: "pending",
    createdAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
  },
  {
    id: "r-5",
    referredName: "Eva Martinez",
    status: "pending",
    createdAt: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000),
  },
  {
    id: "r-6",
    referredName: "Frank Lee",
    status: "rewarded",
    createdAt: new Date(Date.now() - 45 * 24 * 60 * 60 * 1000),
    completedAt: new Date(Date.now() - 40 * 24 * 60 * 60 * 1000),
  },
  {
    id: "r-7",
    referredName: "Grace Chen",
    status: "pending",
    createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000),
  },
];

const PAGE_SIZE = 5;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function relativeTime(date: Date): string {
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

// ---------------------------------------------------------------------------
// Icons
// ---------------------------------------------------------------------------

function ArrowLeftIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
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
  );
}

function CopyIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9.75a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184"
      />
    </svg>
  );
}

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M4.5 12.75l6 6 9-13.5"
      />
    </svg>
  );
}

function ChevronLeftIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M15.75 19.5L8.25 12l7.5-7.5"
      />
    </svg>
  );
}

function ChevronRightIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M8.25 4.5l7.5 7.5-7.5 7.5"
      />
    </svg>
  );
}

function UsersIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"
      />
    </svg>
  );
}

function CheckCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

function GiftIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M21 11.25v8.25a1.5 1.5 0 01-1.5 1.5H5.25a1.5 1.5 0 01-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 109.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1114.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"
      />
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

function StatusBadge({ status }: { status: ReferralStatus }) {
  const styles = {
    pending: "bg-warning/10 text-warning",
    completed: "bg-success/10 text-success",
    rewarded: "bg-primary/10 text-primary",
  };

  const labels = {
    pending: "Pending",
    completed: "Completed",
    rewarded: "Rewarded",
  };

  return (
    <span
      className={`inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium ${styles[status]}`}
    >
      {labels[status]}
    </span>
  );
}

function CopyButton({ text, label }: { text: string; label?: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(text).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <button
      type="button"
      onClick={handleCopy}
      className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-primary text-primary-foreground text-sm font-medium hover:opacity-90 transition-opacity"
    >
      {copied ? (
        <>
          <CheckIcon className="h-4 w-4" />
          Copied
        </>
      ) : (
        <>
          <CopyIcon className="h-4 w-4" />
          {label ?? "Copy"}
        </>
      )}
    </button>
  );
}

function StatCard({
  label,
  value,
  icon,
}: {
  label: string;
  value: number;
  icon: React.ReactNode;
}) {
  return (
    <div className="bg-card border border-border rounded-xl p-4 shadow-sm">
      <div className="flex items-center gap-3">
        <div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
          {icon}
        </div>
        <div>
          <p className="text-2xl font-bold text-foreground">{value}</p>
          <p className="text-sm text-muted-foreground">{label}</p>
        </div>
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

export function ReferralDemo() {
  const [page, setPage] = useState(0);

  const stats = useMemo(() => {
    const total = MOCK_REFERRALS.length;
    const completed = MOCK_REFERRALS.filter(
      (r) => r.status === "completed" || r.status === "rewarded"
    ).length;
    const rewarded = MOCK_REFERRALS.filter(
      (r) => r.status === "rewarded"
    ).length;
    return { total, completed, rewarded };
  }, []);

  const totalPages = Math.max(1, Math.ceil(MOCK_REFERRALS.length / PAGE_SIZE));
  const paginated = MOCK_REFERRALS.slice(
    page * PAGE_SIZE,
    (page + 1) * PAGE_SIZE
  );

  const referralLink = `https://app.cleansaas.dev/register?ref=${MOCK_CODE}`;

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 max-w-2xl">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeftIcon className="h-4 w-4" />
        Back to demos
      </Link>

      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold text-foreground">
            Referral Program
          </h1>
          <p className="text-muted-foreground mt-1">
            Invite friends and earn rewards when they sign up.
          </p>
        </div>

        {/* Referral code card */}
        <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
          <h2 className="text-base font-semibold text-foreground mb-4">
            Your referral code
          </h2>
          <div className="flex items-center gap-3">
            <code className="text-lg font-mono font-bold text-primary bg-primary/5 border border-primary/20 rounded-lg px-4 py-2">
              {MOCK_CODE}
            </code>
            <CopyButton text={MOCK_CODE} />
          </div>

          {/* Share link */}
          <div className="mt-4 pt-4 border-t border-border">
            <p className="text-sm text-muted-foreground mb-2">
              Share this link with friends:
            </p>
            <div className="flex items-center gap-3">
              <input
                readOnly
                value={referralLink}
                className="flex-1 text-sm bg-muted/50 border border-border rounded-lg px-3 py-2 text-muted-foreground"
              />
              <CopyButton text={referralLink} label="Copy link" />
            </div>
          </div>
        </div>

        {/* Stats cards */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <StatCard
            label="Total referrals"
            value={stats.total}
            icon={<UsersIcon className="h-5 w-5" />}
          />
          <StatCard
            label="Completed"
            value={stats.completed}
            icon={<CheckCircleIcon className="h-5 w-5" />}
          />
          <StatCard
            label="Rewards earned"
            value={stats.rewarded}
            icon={<GiftIcon className="h-5 w-5" />}
          />
        </div>

        {/* Referrals table */}
        <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden">
          <div className="px-5 py-4 border-b border-border">
            <h2 className="text-base font-semibold text-foreground">
              Your referrals
            </h2>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border bg-muted/30">
                  <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                    Friend
                  </th>
                  <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                    Status
                  </th>
                  <th className="text-left px-5 py-3 font-medium text-muted-foreground">
                    Date
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {paginated.map((ref) => (
                  <tr
                    key={ref.id}
                    className="hover:bg-muted/20 transition-colors"
                  >
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-2.5">
                        <div className="h-8 w-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold flex-shrink-0">
                          {ref.referredName
                            .split(" ")
                            .map((n) => n[0])
                            .join("")}
                        </div>
                        <span className="text-foreground font-medium">
                          {ref.referredName}
                        </span>
                      </div>
                    </td>
                    <td className="px-5 py-3">
                      <StatusBadge status={ref.status} />
                    </td>
                    <td className="px-5 py-3 text-muted-foreground">
                      {relativeTime(ref.createdAt)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {MOCK_REFERRALS.length > PAGE_SIZE && (
            <div className="flex items-center justify-between px-5 py-3 border-t border-border bg-muted/30">
              <button
                type="button"
                onClick={() => setPage((p) => Math.max(0, p - 1))}
                disabled={page === 0}
                className={cn(
                  "inline-flex items-center gap-1 text-sm font-medium transition-colors",
                  page === 0
                    ? "text-muted-foreground/40 cursor-not-allowed"
                    : "text-foreground hover:text-primary"
                )}
              >
                <ChevronLeftIcon className="h-4 w-4" />
                Previous
              </button>

              <span className="text-xs text-muted-foreground">
                Page {page + 1} of {totalPages}
              </span>

              <button
                type="button"
                onClick={() =>
                  setPage((p) => Math.min(totalPages - 1, p + 1))
                }
                disabled={page >= totalPages - 1}
                className={cn(
                  "inline-flex items-center gap-1 text-sm font-medium transition-colors",
                  page >= totalPages - 1
                    ? "text-muted-foreground/40 cursor-not-allowed"
                    : "text-foreground hover:text-primary"
                )}
              >
                Next
                <ChevronRightIcon className="h-4 w-4" />
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Footer note */}
      <p className="text-center text-sm text-muted-foreground mt-8">
        This is a demo with simulated data. In production, referral codes are
        generated server-side and tracked in the database.
      </p>
    </div>
  );
}
