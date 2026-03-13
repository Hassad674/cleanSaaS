"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { AUTH_TOKEN_KEY } from "@/shared/lib/constants";
import { Loading } from "@/shared/components/loading";

export function GuestGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem(AUTH_TOKEN_KEY);
    if (token) {
      router.push("/dashboard");
    } else {
      setChecking(false);
    }
  }, [router]);

  if (checking) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-muted">
        <Loading />
      </div>
    );
  }

  return <>{children}</>;
}
