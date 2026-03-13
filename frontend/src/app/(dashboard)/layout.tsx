import { AuthGuard } from "@/shared/components/auth-guard";
import { DashboardLayout } from "@/shared/components/layouts/dashboard-layout";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <DashboardLayout>{children}</DashboardLayout>
    </AuthGuard>
  );
}
