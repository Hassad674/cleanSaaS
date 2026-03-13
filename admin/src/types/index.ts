export type User = {
  id: string;
  email: string;
  name: string;
  role: string;
  email_verified: boolean;
  created_at: string;
  updated_at: string;
};

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

export type PaginatedResponse<T> = {
  data: T[];
  total: number;
  page: number;
  limit: number;
};
