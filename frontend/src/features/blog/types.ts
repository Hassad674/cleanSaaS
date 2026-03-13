export type BlogPost = {
  id: string;
  author_id: string;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  cover_image_url: string;
  meta_title: string;
  meta_description: string;
  tags: string[];
  status: string;
  published_at: string | null;
  created_at: string;
  updated_at: string;
};

export type BlogPostsResponse = {
  posts: BlogPost[];
  total: number;
  page: number;
  limit: number;
};

export type BlogTagsResponse = {
  tags: Record<string, number>;
};
