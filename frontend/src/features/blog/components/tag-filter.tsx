"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { cn } from "@/shared/lib/utils";

type TagFilterProps = {
  tags: Record<string, number>;
};

export function TagFilter({ tags }: TagFilterProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const activeTag = searchParams.get("tag") || "";

  const sortedTags = Object.entries(tags).sort(
    ([, a], [, b]) => b - a
  );

  function handleTagClick(tag: string) {
    const params = new URLSearchParams(searchParams.toString());

    if (tag === activeTag) {
      params.delete("tag");
    } else {
      params.set("tag", tag);
    }

    // Reset to page 1 when changing tags
    params.delete("page");

    const query = params.toString();
    router.push(query ? `/blog?${query}` : "/blog");
  }

  if (sortedTags.length === 0) {
    return null;
  }

  return (
    <div className="flex flex-wrap gap-2">
      {/* All posts button */}
      <button
        onClick={() => handleTagClick(activeTag)}
        className={cn(
          "text-sm font-medium rounded-lg px-3 py-1.5 transition-colors",
          !activeTag
            ? "bg-primary text-primary-foreground"
            : "bg-muted text-muted-foreground hover:text-foreground"
        )}
      >
        All
      </button>

      {sortedTags.map(([tag, count]) => (
        <button
          key={tag}
          onClick={() => handleTagClick(tag)}
          className={cn(
            "text-sm font-medium rounded-lg px-3 py-1.5 transition-colors",
            activeTag === tag
              ? "bg-primary text-primary-foreground"
              : "bg-muted text-muted-foreground hover:text-foreground"
          )}
        >
          {tag}
          <span className="ml-1.5 opacity-70">({count})</span>
        </button>
      ))}
    </div>
  );
}
