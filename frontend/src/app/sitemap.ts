import type { MetadataRoute } from "next";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const baseUrl =
    process.env.NEXT_PUBLIC_SITE_URL || "https://cleansaas.dev";

  // Static pages
  const staticPages = ["", "/pricing", "/blog"].map((path) => ({
    url: `${baseUrl}${path}`,
    lastModified: new Date(),
    changeFrequency: "weekly" as const,
    priority: path === "" ? 1 : 0.8,
  }));

  // Dynamic blog posts
  try {
    const apiUrl =
      process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";
    const res = await fetch(`${apiUrl}/blog/posts?limit=100`, {
      next: { revalidate: 3600 },
    });
    const data = await res.json();
    const blogPages = (
      data.posts as Array<{ slug: string; updated_at: string }>
    ).map((post) => ({
      url: `${baseUrl}/blog/${post.slug}`,
      lastModified: new Date(post.updated_at),
      changeFrequency: "monthly" as const,
      priority: 0.6,
    }));
    return [...staticPages, ...blogPages];
  } catch {
    return staticPages;
  }
}
