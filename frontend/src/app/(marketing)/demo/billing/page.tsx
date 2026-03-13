import type { Metadata } from "next";
import { BillingDemo } from "./billing-demo";

export const metadata: Metadata = {
  title: "Billing Demo — CleanSaaS",
  description:
    "Interactive demo of CleanSaaS billing: pricing plans, checkout, and invoices. No account required.",
};

export default function BillingDemoPage() {
  return <BillingDemo />;
}
