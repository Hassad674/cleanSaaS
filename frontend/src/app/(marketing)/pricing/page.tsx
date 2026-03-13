import type { Metadata } from "next";
import { PricingCards } from "@/features/billing/components/pricing-cards";

export const metadata: Metadata = {
  title: "Pricing",
  description: "Simple, transparent pricing for every stage of your business.",
  openGraph: {
    title: "Pricing — CleanSaaS",
    description: "Simple, transparent pricing for every stage of your business.",
    type: "website",
  },
};

export default function PricingPage() {
  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold text-foreground mb-4">Pricing</h1>
        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
          Simple, transparent pricing for every stage of your business.
        </p>
      </div>
      <PricingCards />
    </div>
  );
}
