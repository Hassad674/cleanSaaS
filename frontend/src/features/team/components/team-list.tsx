"use client";

import { useState } from "react";
import { useTeam } from "@/features/team/hooks/use-team";
import { TeamCard } from "@/features/team/components/team-card";
import { CreateTeamDialog } from "@/features/team/components/create-team-dialog";

function LoadingSkeleton() {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          className="bg-card border border-border rounded-xl p-5 shadow-sm animate-pulse"
        >
          <div className="flex items-start gap-4">
            <div className="h-12 w-12 rounded-lg bg-muted flex-shrink-0" />
            <div className="flex-1 space-y-2">
              <div className="h-5 bg-muted rounded w-3/4" />
              <div className="h-4 bg-muted rounded w-1/2" />
            </div>
          </div>
          <div className="flex items-center justify-between mt-4 pt-3 border-t border-border">
            <div className="h-5 bg-muted rounded w-16" />
            <div className="h-4 bg-muted rounded w-20" />
          </div>
        </div>
      ))}
    </div>
  );
}

export function TeamList() {
  const { teams, loading, error, createTeam } = useTeam();
  const [dialogOpen, setDialogOpen] = useState(false);

  if (loading) {
    return <LoadingSkeleton />;
  }

  return (
    <div className="space-y-6">
      {/* Error */}
      {error && (
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3">
          <p className="text-sm text-destructive">{error}</p>
        </div>
      )}

      {/* Header with create button */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Teams</h1>
          <p className="text-muted-foreground mt-1">
            Manage your teams and organizations.
          </p>
        </div>
        <button
          type="button"
          onClick={() => setDialogOpen(true)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground text-sm font-medium rounded-lg hover:opacity-90 transition-opacity"
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
              d="M12 4.5v15m7.5-7.5h-15"
            />
          </svg>
          New team
        </button>
      </div>

      {/* Teams grid */}
      {teams.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 px-4">
          <div className="h-16 w-16 rounded-xl bg-muted flex items-center justify-center mb-4">
            <svg
              className="h-8 w-8 text-muted-foreground/50"
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
          </div>
          <h3 className="text-base font-medium text-foreground mb-1">
            No teams yet
          </h3>
          <p className="text-sm text-muted-foreground text-center mb-4">
            Create a team to collaborate with others.
          </p>
          <button
            type="button"
            onClick={() => setDialogOpen(true)}
            className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground text-sm font-medium rounded-lg hover:opacity-90 transition-opacity"
          >
            Create your first team
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {teams.map((team) => (
            <TeamCard key={team.id} team={team} />
          ))}
        </div>
      )}

      {/* Create dialog */}
      <CreateTeamDialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        onSubmit={createTeam}
      />
    </div>
  );
}
