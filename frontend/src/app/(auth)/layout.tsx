import { GuestGuard } from "@/shared/components/guest-guard";
import { AuthLayout } from "@/shared/components/layouts/auth-layout";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <GuestGuard>
      <AuthLayout>{children}</AuthLayout>
    </GuestGuard>
  );
}
