const SITE_URL =
  process.env.NEXT_PUBLIC_SITE_URL || "https://cleansaas.dev";

type OrganizationJsonLd = {
  "@context": "https://schema.org";
  "@type": "Organization";
  name: string;
  url: string;
  description: string;
  logo?: string;
};

export function organizationJsonLd(): OrganizationJsonLd {
  return {
    "@context": "https://schema.org",
    "@type": "Organization",
    name: "CleanSaaS",
    url: SITE_URL,
    description:
      "Open-source SaaS boilerplate for modern applications. Ship your SaaS in weeks, not months.",
  };
}
