"use client";

import { useState } from "react";
import { useAuth } from "@/shared/hooks/use-auth";
import { deleteAccount } from "@/features/user/actions/user";

export function SettingsDanger() {
  const { getToken, logout } = useAuth({ required: true });

  const [showConfirm, setShowConfirm] = useState(false);
  const [confirmText, setConfirmText] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleDelete() {
    if (confirmText !== "DELETE") return;

    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);

    const res = await deleteAccount(token);
    setLoading(false);

    if (res.error) {
      setError(res.error);
    } else {
      logout();
    }
  }

  return (
    <div className="bg-card border border-destructive rounded-xl p-6 shadow-sm">
      <h2 className="text-lg font-semibold text-foreground mb-4">
        Danger Zone
      </h2>
      <p className="text-sm text-muted-foreground mb-4">
        Permanently delete your account and all associated data. This action
        cannot be undone.
      </p>

      {!showConfirm ? (
        <button
          onClick={() => setShowConfirm(true)}
          className="bg-destructive text-destructive-foreground rounded-lg px-4 py-2 hover:opacity-90 transition-opacity"
        >
          Delete my account
        </button>
      ) : (
        <div className="space-y-3">
          <div>
            <label
              htmlFor="confirm-delete"
              className="block text-sm font-medium text-foreground mb-1"
            >
              Type <span className="font-bold">DELETE</span> to confirm
            </label>
            <input
              id="confirm-delete"
              type="text"
              value={confirmText}
              onChange={(e) => setConfirmText(e.target.value)}
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              placeholder="DELETE"
            />
          </div>

          {error && <p className="text-sm text-destructive">{error}</p>}

          <div className="flex gap-3">
            <button
              onClick={handleDelete}
              disabled={confirmText !== "DELETE" || loading}
              className="bg-destructive text-destructive-foreground rounded-lg px-4 py-2 hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              {loading ? "Deleting..." : "Confirm deletion"}
            </button>
            <button
              onClick={() => {
                setShowConfirm(false);
                setConfirmText("");
                setError(null);
              }}
              className="border border-border text-foreground rounded-lg px-4 py-2 hover:opacity-90 transition-opacity"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
