import Link from "next/link";
import { formatDate } from "@/shared/lib/utils";
import type { BlogPost } from "@/features/blog/types";

type PostCardProps = {
  post: BlogPost;
};

export function PostCard({ post }: PostCardProps) {
  return (
    <Link
      href={`/blog/${post.slug}`}
      className="group bg-card border border-border rounded-xl shadow-sm overflow-hidden flex flex-col hover:shadow-md transition-shadow"
    >
      {/* Cover image */}
      {post.cover_image_url && (
        <div className="h-48 w-full overflow-hidden">
          <img
            src={post.cover_image_url}
            alt={post.title}
            className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-300"
            loading="lazy"
          />
        </div>
      )}

      {/* Content */}
      <div className="p-6 flex-1 flex flex-col gap-3">
        <h3 className="text-lg font-bold text-foreground group-hover:text-primary transition-colors line-clamp-2">
          {post.title}
        </h3>

        {post.excerpt && (
          <p className="text-sm text-muted-foreground line-clamp-3">
            {post.excerpt}
          </p>
        )}

        <div className="mt-auto pt-3 flex flex-wrap items-center gap-2">
          {/* Tags */}
          {post.tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5">
              {post.tags.slice(0, 3).map((tag) => (
                <span
                  key={tag}
                  className="text-xs font-medium bg-accent text-foreground rounded-md px-2 py-0.5"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}

          {/* Date */}
          {post.published_at && (
            <span className="text-xs text-muted-foreground ml-auto">
              {formatDate(post.published_at)}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}
