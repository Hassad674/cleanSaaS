import type { Metadata } from "next";
import { ReferralDashboard } from "@/features/referral/components/referral-dashboard";

export const metadata: Metadata = {
  title: "Referral Program",
  robots: { index: false, follow: false },
};

export default function ReferralPage() {
  return (
    <div className="max-w-2xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Referral Program</h1>
        <p className="text-muted-foreground mt-1">
          Invite friends and earn rewards when they sign up.
        </p>
      </div>
      <ReferralDashboard />
    </div>
  );
}
