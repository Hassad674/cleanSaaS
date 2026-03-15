import type { Metadata } from "next";
import { Suspense } from "react";
import { BillingDemo } from "./billing-demo";

export const metadata: Metadata = {
  title: "Billing Demo — CleanSaaS",
  description:
    "Live Stripe Checkout demo with monthly, yearly, and lifetime plans. No account required.",
};

export default function BillingDemoPage() {
  return (
    <Suspense
      fallback={
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <div className="animate-pulse h-96" />
        </div>
      }
    >
      <BillingDemo />
    </Suspense>
  );
}
