"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { verifyEmail } from "@/features/auth/actions/auth";

export function VerifyEmailForm() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token") ?? "";

  const [status, setStatus] = useState<"loading" | "success" | "error" | "no-token">(
    token ? "loading" : "no-token"
  );
  const [error, setError] = useState("");

  useEffect(() => {
    if (!token) return;

    async function verify() {
      const res = await verifyEmail(token);
      if (res.error) {
        setError(res.error);
        setStatus("error");
      } else {
        setStatus("success");
      }
    }

    verify();
  }, [token]);

  return (
    <div className="space-y-4 text-center">
      <h1 className="text-2xl font-bold text-card-foreground">
        Email verification
      </h1>

      {status === "loading" && (
        <p className="text-muted-foreground">Verifying your email...</p>
      )}

      {status === "success" && (
        <div className="space-y-3">
          <div className="mx-auto w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
            <svg className="w-6 h-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <p className="text-foreground font-medium">
            Your email has been verified successfully!
          </p>
          <Link
            href="/dashboard"
            className="inline-block bg-primary text-primary-foreground px-6 py-2 rounded-lg hover:opacity-90 transition-opacity font-medium"
          >
            Go to dashboard
          </Link>
        </div>
      )}

      {status === "error" && (
        <div className="space-y-3">
          <p className="text-destructive">{error}</p>
          <p className="text-sm text-muted-foreground">
            The verification link may have expired or already been used.
          </p>
          <Link
            href="/login"
            className="text-sm text-primary hover:underline"
          >
            Back to login
          </Link>
        </div>
      )}

      {status === "no-token" && (
        <div className="space-y-3">
          <p className="text-muted-foreground">
            No verification token found. Please check your email for the verification link.
          </p>
          <Link
            href="/login"
            className="text-sm text-primary hover:underline"
          >
            Back to login
          </Link>
        </div>
      )}
    </div>
  );
}
