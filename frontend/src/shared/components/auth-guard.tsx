"use client";

import { useAuth } from "@/features/auth/hooks/use-auth";
import { Loading } from "@/shared/components/loading";

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth({ required: true });

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Loading />
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return <>{children}</>;
}
