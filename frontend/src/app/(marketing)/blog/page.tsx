import type { Metadata } from "next";
import { fetchPublishedPosts, fetchTags } from "@/features/blog/actions/blog";
import { BlogList } from "@/features/blog/components/blog-list";

export const metadata: Metadata = {
  title: "Blog",
  description: "Latest articles, tutorials, and updates from our team.",
};

type BlogPageProps = {
  searchParams: Promise<{ tag?: string; page?: string }>;
};

export default async function BlogPage({ searchParams }: BlogPageProps) {
  const params = await searchParams;
  const tag = params.tag;
  const page = Number(params.page) || 1;
  const limit = 12;

  const [postsRes, tagsRes] = await Promise.all([
    fetchPublishedPosts(tag, page, limit),
    fetchTags(),
  ]);

  const posts = postsRes.data?.posts ?? [];
  const total = postsRes.data?.total ?? 0;
  const tags = tagsRes.data?.tags ?? {};

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-foreground mb-2">Blog</h1>
        <p className="text-lg text-muted-foreground">
          Latest articles, tutorials, and updates from our team.
        </p>
      </div>
      <BlogList
        posts={posts}
        tags={tags}
        total={total}
        page={page}
        limit={limit}
        activeTag={tag}
      />
    </div>
  );
}
