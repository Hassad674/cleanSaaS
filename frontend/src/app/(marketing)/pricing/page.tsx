import type { Metadata } from "next";

export const metadata: Metadata = { title: "Pricing" };

export default function PricingPage() {
  return (
    <div className="container mx-auto px-4 py-16 text-center">
      <h1 className="text-4xl font-bold mb-4">Pricing</h1>
      <p className="text-muted-foreground">Plans will be configured here.</p>
    </div>
  );
}
