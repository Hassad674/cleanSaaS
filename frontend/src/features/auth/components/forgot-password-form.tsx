"use client";

import { useState } from "react";
import Link from "next/link";
import { forgotPassword } from "@/features/auth/actions/auth";

export function ForgotPasswordForm() {
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setSuccess(false);
    setLoading(true);

    const res = await forgotPassword(email);

    if (res.error) {
      setError(res.error);
      setLoading(false);
      return;
    }

    setSuccess(true);
    setLoading(false);
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h1 className="text-2xl font-bold text-center text-card-foreground">
        Forgot password
      </h1>

      <p className="text-sm text-muted-foreground text-center">
        Enter your email and we&apos;ll send you a link to reset your password.
      </p>

      {error && (
        <p className="text-sm text-destructive text-center">{error}</p>
      )}

      {success && (
        <p className="text-sm text-green-600 text-center">
          If an account exists with that email, we&apos;ve sent a reset link.
        </p>
      )}

      <div>
        <label
          htmlFor="email"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full px-3 py-2 border border-border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-primary text-primary-foreground py-2 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 font-medium"
      >
        {loading ? "Sending..." : "Send reset link"}
      </button>

      <div className="text-center text-sm">
        <Link href="/login" className="text-primary hover:underline">
          Back to login
        </Link>
      </div>
    </form>
  );
}
