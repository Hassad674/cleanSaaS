import type { BlogPost } from "@/features/blog/types";

const SITE_URL =
  process.env.NEXT_PUBLIC_SITE_URL || "https://cleansaas.dev";

type ArticleJsonLd = {
  "@context": "https://schema.org";
  "@type": "Article";
  headline: string;
  description: string;
  url: string;
  datePublished?: string;
  dateModified: string;
  image?: string;
  author: {
    "@type": "Organization";
    name: string;
  };
  publisher: {
    "@type": "Organization";
    name: string;
  };
};

export function articleJsonLd(post: BlogPost): ArticleJsonLd {
  return {
    "@context": "https://schema.org",
    "@type": "Article",
    headline: post.meta_title || post.title,
    description: post.meta_description || post.excerpt,
    url: `${SITE_URL}/blog/${post.slug}`,
    ...(post.published_at && { datePublished: post.published_at }),
    dateModified: post.updated_at,
    ...(post.cover_image_url && { image: post.cover_image_url }),
    author: {
      "@type": "Organization",
      name: "CleanSaaS",
    },
    publisher: {
      "@type": "Organization",
      name: "CleanSaaS",
    },
  };
}
