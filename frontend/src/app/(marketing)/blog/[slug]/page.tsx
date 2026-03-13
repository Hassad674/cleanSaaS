import type { Metadata } from "next";
import { notFound } from "next/navigation";
import Link from "next/link";
import { fetchPostBySlug } from "@/features/blog/actions/blog";
import { PostContent } from "@/features/blog/components/post-content";
import { articleJsonLd } from "@/features/blog/lib/json-ld";

type PostPageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateMetadata({
  params,
}: PostPageProps): Promise<Metadata> {
  const { slug } = await params;
  const { data: post } = await fetchPostBySlug(slug);

  if (!post) {
    return { title: "Post Not Found" };
  }

  return {
    title: post.meta_title || post.title,
    description: post.meta_description || post.excerpt,
    openGraph: {
      title: post.meta_title || post.title,
      description: post.meta_description || post.excerpt,
      type: "article",
      publishedTime: post.published_at ?? undefined,
      ...(post.cover_image_url && {
        images: [{ url: post.cover_image_url }],
      }),
    },
    twitter: {
      card: "summary_large_image",
      title: post.meta_title || post.title,
      description: post.meta_description || post.excerpt,
      ...(post.cover_image_url && {
        images: [post.cover_image_url],
      }),
    },
  };
}

export default async function PostPage({ params }: PostPageProps) {
  const { slug } = await params;
  const { data: post, error } = await fetchPostBySlug(slug);

  if (!post || error) {
    notFound();
  }

  const jsonLd = articleJsonLd(post);

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />
      <nav className="mb-8">
        <Link
          href="/blog"
          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          &larr; Back to blog
        </Link>
      </nav>
      <PostContent post={post} />
    </div>
  );
}
