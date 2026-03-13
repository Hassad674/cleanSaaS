import { Hero } from "@/features/marketing/components/hero";
import { StackSection } from "@/features/marketing/components/stack-section";
import { FeaturesSection } from "@/features/marketing/components/features-section";
import { ArchitectureSection } from "@/features/marketing/components/architecture-section";
import { DXSection } from "@/features/marketing/components/dx-section";
import { CTASection } from "@/features/marketing/components/cta-section";

export default function HomePage() {
  return (
    <>
      <Hero />
      <StackSection />
      <FeaturesSection />
      <ArchitectureSection />
      <DXSection />
      <CTASection />
    </>
  );
}
