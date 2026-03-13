import Link from "next/link";
import { cn } from "@/shared/lib/utils";
import { PostCard } from "@/features/blog/components/post-card";
import { TagFilter } from "@/features/blog/components/tag-filter";
import type { BlogPost } from "@/features/blog/types";

type BlogListProps = {
  posts: BlogPost[];
  tags: Record<string, number>;
  total: number;
  page: number;
  limit: number;
  activeTag?: string;
};

export function BlogList({
  posts,
  tags,
  total,
  page,
  limit,
  activeTag,
}: BlogListProps) {
  const totalPages = Math.ceil(total / limit);
  const hasPrev = page > 1;
  const hasNext = page < totalPages;

  function buildPageUrl(targetPage: number): string {
    const params = new URLSearchParams();
    if (activeTag) params.set("tag", activeTag);
    if (targetPage > 1) params.set("page", String(targetPage));
    const query = params.toString();
    return query ? `/blog?${query}` : "/blog";
  }

  return (
    <div className="space-y-8">
      {/* Tag filter */}
      <TagFilter tags={tags} />

      {/* Post grid */}
      {posts.length === 0 ? (
        <div className="bg-card border border-border rounded-xl p-12 shadow-sm text-center">
          <svg
            className="mx-auto h-12 w-12 text-muted-foreground mb-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={1.5}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M12 7.5h1.5m-1.5 3h1.5m-7.5 3h7.5m-7.5 3h7.5m3-9h3.375c.621 0 1.125.504 1.125 1.125V18a2.25 2.25 0 01-2.25 2.25M16.5 7.5V18a2.25 2.25 0 002.25 2.25M16.5 7.5V4.875c0-.621-.504-1.125-1.125-1.125H4.125C3.504 3.75 3 4.254 3 4.875V18a2.25 2.25 0 002.25 2.25h13.5M6 7.5h3v3H6V7.5z"
            />
          </svg>
          <h3 className="text-base font-medium text-foreground mb-1">
            No posts yet
          </h3>
          <p className="text-sm text-muted-foreground">
            {activeTag
              ? `No posts found with tag "${activeTag}".`
              : "Blog posts will appear here once published."}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between pt-2">
          {hasPrev ? (
            <Link
              href={buildPageUrl(page - 1)}
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Previous
            </Link>
          ) : (
            <span className="text-sm text-muted-foreground opacity-50">
              Previous
            </span>
          )}

          <span className="text-xs text-muted-foreground">
            Page {page} of {totalPages} &middot; {total} post
            {total !== 1 ? "s" : ""}
          </span>

          {hasNext ? (
            <Link
              href={buildPageUrl(page + 1)}
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Next
            </Link>
          ) : (
            <span className="text-sm text-muted-foreground opacity-50">
              Next
            </span>
          )}
        </div>
      )}
    </div>
  );
}
