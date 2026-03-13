"use server";

import { api } from "@/shared/lib/api";
import type { BlogPost, BlogPostsResponse, BlogTagsResponse } from "@/features/blog/types";

export async function fetchPublishedPosts(
  tag?: string,
  page: number = 1,
  limit: number = 12
) {
  const params = new URLSearchParams();
  if (tag) params.set("tag", tag);
  params.set("page", String(page));
  params.set("limit", String(limit));

  return api<BlogPostsResponse>(`/blog/posts?${params.toString()}`);
}

export async function fetchPostBySlug(slug: string) {
  return api<BlogPost>(`/blog/posts/${encodeURIComponent(slug)}`);
}

export async function fetchTags() {
  return api<BlogTagsResponse>("/blog/tags");
}
