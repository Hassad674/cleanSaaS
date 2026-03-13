import type { MetadataRoute } from "next";

export default function robots(): MetadataRoute.Robots {
  const baseUrl =
    process.env.NEXT_PUBLIC_SITE_URL || "https://cleansaas.dev";

  return {
    rules: [
      {
        userAgent: "*",
        allow: ["/", "/blog", "/pricing"],
        disallow: ["/dashboard", "/settings", "/admin", "/api", "/auth"],
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
  };
}
