import { AuthLayout } from "@/shared/components/layouts/auth-layout";

export default function Layout({ children }: { children: React.ReactNode }) {
  return <AuthLayout>{children}</AuthLayout>;
}
