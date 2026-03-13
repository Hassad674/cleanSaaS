"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import type { Team } from "@/features/team/types";

export function TeamSettings({
  team,
  isOwner,
  onUpdate,
  onDelete,
  onLeave,
}: {
  team: Team;
  isOwner: boolean;
  onUpdate: (name: string) => Promise<{ success: boolean; error: string | null }>;
  onDelete: () => Promise<{ success: boolean; error: string | null }>;
  onLeave: () => Promise<{ success: boolean; error: string | null }>;
}) {
  const [name, setName] = useState(team.name);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [confirmLeave, setConfirmLeave] = useState(false);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed || trimmed === team.name) return;

    setSaving(true);
    setError(null);
    setSuccess(false);

    const result = await onUpdate(trimmed);
    if (result.success) {
      setSuccess(true);
      setTimeout(() => setSuccess(false), 3000);
    } else {
      setError(result.error);
    }
    setSaving(false);
  };

  const handleDelete = async () => {
    const result = await onDelete();
    if (!result.success) {
      setError(result.error);
    }
  };

  const handleLeave = async () => {
    const result = await onLeave();
    if (!result.success) {
      setError(result.error);
    }
  };

  return (
    <div className="space-y-6">
      {/* Error */}
      {error && (
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3">
          <p className="text-sm text-destructive">{error}</p>
        </div>
      )}

      {/* Success */}
      {success && (
        <div className="bg-success/10 border border-success/20 rounded-lg px-4 py-3">
          <p className="text-sm text-success">Team updated successfully.</p>
        </div>
      )}

      {/* Team name */}
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
        <h2 className="text-base font-semibold text-foreground mb-4">
          Team settings
        </h2>
        <form onSubmit={handleSave}>
          <label
            htmlFor="team-name-edit"
            className="block text-sm font-medium text-foreground mb-1.5"
          >
            Team name
          </label>
          <input
            id="team-name-edit"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            disabled={!isOwner}
            className={cn(
              "w-full bg-muted/50 border border-border rounded-lg px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
              !isOwner && "opacity-60 cursor-not-allowed"
            )}
          />

          {/* Team info */}
          <div className="flex items-center gap-4 mt-3 text-xs text-muted-foreground">
            <span>Slug: {team.slug}</span>
            <span>Plan: {team.plan}</span>
            <span>Max members: {team.max_members}</span>
          </div>

          {isOwner && (
            <div className="mt-4">
              <button
                type="submit"
                disabled={!name.trim() || name.trim() === team.name || saving}
                className={cn(
                  "px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity",
                  (!name.trim() || name.trim() === team.name || saving) &&
                    "opacity-50 cursor-not-allowed"
                )}
              >
                {saving ? "Saving..." : "Save changes"}
              </button>
            </div>
          )}
        </form>
      </div>

      {/* Danger zone */}
      <div className="bg-card border border-destructive/30 rounded-xl p-6 shadow-sm">
        <h2 className="text-base font-semibold text-destructive mb-2">
          Danger zone
        </h2>

        {isOwner ? (
          <div>
            <p className="text-sm text-muted-foreground mb-4">
              Permanently delete this team and all associated data. This action
              cannot be undone.
            </p>
            {confirmDelete ? (
              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={handleDelete}
                  className="px-4 py-2 text-sm font-medium rounded-lg bg-destructive text-destructive-foreground hover:opacity-90 transition-opacity"
                >
                  Yes, delete team
                </button>
                <button
                  type="button"
                  onClick={() => setConfirmDelete(false)}
                  className="px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors rounded-lg"
                >
                  Cancel
                </button>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setConfirmDelete(true)}
                className="px-4 py-2 text-sm font-medium rounded-lg border border-destructive text-destructive hover:bg-destructive/10 transition-colors"
              >
                Delete team
              </button>
            )}
          </div>
        ) : (
          <div>
            <p className="text-sm text-muted-foreground mb-4">
              Leave this team. You will lose access to all team resources.
            </p>
            {confirmLeave ? (
              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={handleLeave}
                  className="px-4 py-2 text-sm font-medium rounded-lg bg-destructive text-destructive-foreground hover:opacity-90 transition-opacity"
                >
                  Yes, leave team
                </button>
                <button
                  type="button"
                  onClick={() => setConfirmLeave(false)}
                  className="px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors rounded-lg"
                >
                  Cancel
                </button>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setConfirmLeave(true)}
                className="px-4 py-2 text-sm font-medium rounded-lg border border-destructive text-destructive hover:bg-destructive/10 transition-colors"
              >
                Leave team
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
