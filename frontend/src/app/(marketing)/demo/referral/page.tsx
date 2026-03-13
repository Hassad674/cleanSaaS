import type { Metadata } from "next";
import { ReferralDemo } from "./referral-demo";

export const metadata: Metadata = {
  title: "Referral Program Demo — CleanSaaS",
  description:
    "Interactive referral program demo with code sharing, stats dashboard, and referral tracking. No account required.",
};

export default function ReferralDemoPage() {
  return <ReferralDemo />;
}
