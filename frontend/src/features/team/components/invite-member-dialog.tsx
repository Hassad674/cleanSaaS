"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";

export function InviteMemberDialog({
  open,
  onClose,
  onSubmit,
}: {
  open: boolean;
  onClose: () => void;
  onSubmit: (
    email: string,
    role: "admin" | "member"
  ) => Promise<{ success: boolean; error: string | null }>;
}) {
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<"admin" | "member">("member");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!open) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = email.trim();
    if (!trimmed) return;

    setSubmitting(true);
    setError(null);

    const result = await onSubmit(trimmed, role);
    if (result.success) {
      setEmail("");
      setRole("member");
      onClose();
    } else {
      setError(result.error);
    }
    setSubmitting(false);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Overlay */}
      <div
        className="absolute inset-0 bg-background/80 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Dialog */}
      <div className="relative bg-card border border-border rounded-xl shadow-lg w-full max-w-md mx-4 p-6">
        <h2 className="text-lg font-semibold text-foreground mb-1">
          Invite a team member
        </h2>
        <p className="text-sm text-muted-foreground mb-6">
          Send an invitation by email. They&apos;ll receive a link to join your
          team.
        </p>

        {error && (
          <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3 mb-4">
            <p className="text-sm text-destructive">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          {/* Email input */}
          <label
            htmlFor="invite-email"
            className="block text-sm font-medium text-foreground mb-1.5"
          >
            Email address
          </label>
          <input
            id="invite-email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
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
                  onClick={() => setRole(r)}
                  className={cn(
                    "flex-1 px-3 py-2 text-sm font-medium rounded-lg border transition-colors",
                    role === r
                      ? "border-primary bg-primary/10 text-primary"
                      : "border-border bg-muted/30 text-muted-foreground hover:text-foreground"
                  )}
                >
                  {r === "admin" ? "Admin" : "Member"}
                </button>
              ))}
            </div>
            <p className="text-xs text-muted-foreground mt-2">
              {role === "admin"
                ? "Admins can manage members and team settings."
                : "Members can view and collaborate within the team."}
            </p>
          </fieldset>

          <div className="flex items-center justify-end gap-3 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors rounded-lg"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={!email.trim() || submitting}
              className={cn(
                "px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity",
                (!email.trim() || submitting) &&
                  "opacity-50 cursor-not-allowed"
              )}
            >
              {submitting ? "Sending..." : "Send invite"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
