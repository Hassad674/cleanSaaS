import type { Metadata } from "next";
import { Hero } from "@/features/marketing/components/hero";
import { SpotlightAI } from "@/features/marketing/components/spotlight-ai";
import { SpotlightAdmin } from "@/features/marketing/components/spotlight-admin";
import { SpotlightArchitecture } from "@/features/marketing/components/spotlight-architecture";
import { StackSection } from "@/features/marketing/components/stack-section";
import { ComparisonSection } from "@/features/marketing/components/comparison-section";
import { FeaturesSection } from "@/features/marketing/components/features-section";
import { ArchitectureSection } from "@/features/marketing/components/architecture-section";
import { DXSection } from "@/features/marketing/components/dx-section";
import { CTASection } from "@/features/marketing/components/cta-section";
import { organizationJsonLd } from "@/shared/lib/json-ld";

export const metadata: Metadata = {
  title: "Ship your SaaS in weeks, not months",
  description:
    "Open-source SaaS boilerplate with auth, billing, AI, admin, notifications, and storage. Built with Next.js, Go, and PostgreSQL.",
  openGraph: {
    title: "CleanSaaS — Ship your SaaS in weeks, not months",
    description:
      "Open-source SaaS boilerplate with auth, billing, AI, admin, notifications, and storage. Built with Next.js, Go, and PostgreSQL.",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "CleanSaaS — Ship your SaaS in weeks, not months",
    description:
      "Open-source SaaS boilerplate with auth, billing, AI, admin, notifications, and storage. Built with Next.js, Go, and PostgreSQL.",
  },
};

export default function HomePage() {
  const jsonLd = organizationJsonLd();

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />
      <Hero />
      <FeaturesSection />
      <SpotlightAI />
      <SpotlightAdmin />
      <SpotlightArchitecture />
      <StackSection />
      <ComparisonSection />
      <DXSection />
      <CTASection />
    </>
  );
}
