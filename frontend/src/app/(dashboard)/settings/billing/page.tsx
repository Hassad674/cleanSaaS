import type { Metadata } from "next";
import { SubscriptionStatus } from "@/features/billing/components/subscription-status";
import { InvoiceList } from "@/features/billing/components/invoice-list";

export const metadata: Metadata = {
  title: "Billing",
  robots: { index: false, follow: false },
};

export default function BillingPage() {
  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold text-foreground">Billing</h1>
      <SubscriptionStatus />
      <InvoiceList />
    </div>
  );
}
