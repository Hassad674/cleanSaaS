"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";

export function CreateTeamDialog({
  open,
  onClose,
  onSubmit,
}: {
  open: boolean;
  onClose: () => void;
  onSubmit: (name: string) => Promise<{ success: boolean; error: string | null }>;
}) {
  const [name, setName] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!open) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) return;

    setSubmitting(true);
    setError(null);

    const result = await onSubmit(trimmed);
    if (result.success) {
      setName("");
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
          Create a new team
        </h2>
        <p className="text-sm text-muted-foreground mb-6">
          Give your team a name to get started. You can change it later.
        </p>

        {error && (
          <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3 mb-4">
            <p className="text-sm text-destructive">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <label
            htmlFor="team-name"
            className="block text-sm font-medium text-foreground mb-1.5"
          >
            Team name
          </label>
          <input
            id="team-name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. Acme Corp"
            className="w-full bg-muted/50 border border-border rounded-lg px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            autoFocus
          />

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
              disabled={!name.trim() || submitting}
              className={cn(
                "px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity",
                (!name.trim() || submitting) && "opacity-50 cursor-not-allowed"
              )}
            >
              {submitting ? "Creating..." : "Create team"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
