"use client";

import { useState, type FormEvent } from "react";
import { useAuth } from "@/shared/hooks/use-auth";
import { updateProfile } from "@/features/user/actions/user";

export function SettingsProfile() {
  const { user, getToken } = useAuth({ required: true });

  const [name, setName] = useState(user?.name ?? "");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setSuccess(null);
    setError(null);

    const token = getToken();
    if (!token) return;

    setLoading(true);
    const res = await updateProfile(token, { name });
    setLoading(false);

    if (res.error) {
      setError(res.error);
    } else {
      setSuccess("Profile updated successfully.");
    }
  }

  if (!user) return null;

  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
      <h2 className="text-lg font-semibold text-foreground mb-4">Profile</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="name"
            className="block text-sm font-medium text-foreground mb-1"
          >
            Name
          </label>
          <input
            id="name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            placeholder="Your name"
          />
        </div>

        {success && <p className="text-sm text-green-600">{success}</p>}
        {error && <p className="text-sm text-destructive">{error}</p>}

        <button
          type="submit"
          disabled={loading}
          className="bg-primary text-primary-foreground rounded-lg px-4 py-2 hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {loading ? "Saving..." : "Save"}
        </button>
      </form>
    </div>
  );
}
