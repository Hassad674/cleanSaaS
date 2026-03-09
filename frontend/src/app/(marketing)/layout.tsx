import { MarketingLayout } from "@/shared/components/layouts/marketing-layout";

export default function Layout({ children }: { children: React.ReactNode }) {
  return <MarketingLayout>{children}</MarketingLayout>;
}
