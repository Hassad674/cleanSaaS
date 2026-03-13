import type { User, BlogPost, PaginatedResponse } from "../types/index.ts";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8081";

function getToken(): string | null {
  return localStorage.getItem("admin_token");
}

async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> | undefined),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`API error ${response.status}: ${body}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

export type LoginResponse = {
  token: string;
  user: User;
};

export async function login(
  email: string,
  password: string
): Promise<LoginResponse> {
  return apiFetch<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function fetchCurrentUser(): Promise<User> {
  return apiFetch<User>("/users/me");
}

// ---------------------------------------------------------------------------
// Users (admin)
// ---------------------------------------------------------------------------

type UsersListResponse = {
  users: User[];
  total: number;
  page: number;
  limit: number;
};

export async function fetchUsers(params: {
  page?: number;
  limit?: number;
  search?: string;
}): Promise<PaginatedResponse<User>> {
  const query = new URLSearchParams();
  if (params.page) query.set("page", String(params.page));
  if (params.limit) query.set("limit", String(params.limit));
  if (params.search) query.set("search", params.search);

  const res = await apiFetch<UsersListResponse>(
    `/admin/users?${query.toString()}`
  );
  return { data: res.users, total: res.total, page: res.page, limit: res.limit };
}

export async function updateUserRole(
  userId: string,
  role: string
): Promise<User> {
  return apiFetch<User>(`/admin/users/${userId}/role`, {
    method: "PUT",
    body: JSON.stringify({ role }),
  });
}

export async function updateUserStatus(
  userId: string,
  status: string
): Promise<User> {
  return apiFetch<User>(`/admin/users/${userId}/status`, {
    method: "PUT",
    body: JSON.stringify({ status }),
  });
}

// ---------------------------------------------------------------------------
// Blog posts (admin)
// ---------------------------------------------------------------------------

type BlogPostsListResponse = {
  posts: BlogPost[];
  total: number;
  page: number;
  limit: number;
};

export async function fetchBlogPosts(params: {
  page?: number;
  limit?: number;
  status?: string;
  tag?: string;
}): Promise<PaginatedResponse<BlogPost>> {
  const query = new URLSearchParams();
  if (params.page) query.set("page", String(params.page));
  if (params.limit) query.set("limit", String(params.limit));
  if (params.status) query.set("status", params.status);
  if (params.tag) query.set("tag", params.tag);

  const res = await apiFetch<BlogPostsListResponse>(
    `/admin/blog/posts?${query.toString()}`
  );
  return {
    data: res.posts,
    total: res.total,
    page: res.page,
    limit: res.limit,
  };
}

export type BlogPostInput = {
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  cover_image_url: string;
  meta_title: string;
  meta_description: string;
  tags: string[];
  status: string;
};

export async function createBlogPost(
  input: BlogPostInput
): Promise<BlogPost> {
  return apiFetch<BlogPost>("/admin/blog/posts", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function updateBlogPost(
  id: string,
  input: BlogPostInput
): Promise<BlogPost> {
  return apiFetch<BlogPost>(`/admin/blog/posts/${id}`, {
    method: "PUT",
    body: JSON.stringify(input),
  });
}

export async function deleteBlogPost(id: string): Promise<void> {
  return apiFetch<void>(`/admin/blog/posts/${id}`, {
    method: "DELETE",
  });
}

// ---------------------------------------------------------------------------
// Blog tags
// ---------------------------------------------------------------------------

type TagsResponse = {
  tags: Record<string, number>;
};

export async function fetchBlogTags(): Promise<Record<string, number>> {
  const res = await apiFetch<TagsResponse>("/blog/tags");
  return res.tags;
}
