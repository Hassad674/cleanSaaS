"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { login } from "@/features/auth/actions/auth";
import { useAuth } from "@/features/auth/hooks/use-auth";

export function LoginForm() {
  const router = useRouter();
  const { setToken } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    const res = await login(email, password);

    if (res.error) {
      setError(res.error);
      setLoading(false);
      return;
    }

    if (res.data) {
      setToken(res.data.token);
      router.push("/");
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h1 className="text-2xl font-bold text-center">Log in</h1>

      {error && (
        <p className="text-sm text-red-500 text-center">{error}</p>
      )}

      <div>
        <label htmlFor="email" className="block text-sm font-medium mb-1">
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full px-3 py-2 border border-zinc-300 rounded-md dark:border-zinc-700 dark:bg-zinc-900"
        />
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium mb-1">
          Password
        </label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full px-3 py-2 border border-zinc-300 rounded-md dark:border-zinc-700 dark:bg-zinc-900"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-zinc-900 text-white py-2 rounded-md hover:bg-zinc-800 disabled:opacity-50 dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
      >
        {loading ? "Logging in..." : "Log in"}
      </button>

      <div className="text-center text-sm text-zinc-500 space-y-2">
        <Link href="/forgot-password" className="hover:underline block">
          Forgot password?
        </Link>
        <p>
          Don&apos;t have an account?{" "}
          <Link href="/register" className="font-medium text-zinc-900 dark:text-zinc-100 hover:underline">
            Sign up
          </Link>
        </p>
      </div>
    </form>
  );
}
