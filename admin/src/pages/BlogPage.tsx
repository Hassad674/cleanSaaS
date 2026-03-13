import { useEffect, useState, useCallback, type FormEvent } from "react";
import DataTable, { type Column } from "../components/DataTable.tsx";
import {
  fetchBlogPosts,
  createBlogPost,
  updateBlogPost,
  deleteBlogPost,
  type BlogPostInput,
} from "../lib/api.ts";
import type { BlogPost } from "../types/index.ts";

const LIMIT = 20;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, "")
    .replace(/[\s_]+/g, "-")
    .replace(/-+/g, "-");
}

function emptyForm(): BlogPostInput {
  return {
    title: "",
    slug: "",
    excerpt: "",
    content: "",
    cover_image_url: "",
    meta_title: "",
    meta_description: "",
    tags: [],
    status: "draft",
  };
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

export default function BlogPage() {
  // List state
  const [posts, setPosts] = useState<BlogPost[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [isLoading, setIsLoading] = useState(true);

  // Form state
  const [view, setView] = useState<"list" | "form">("list");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<BlogPostInput>(emptyForm());
  const [tagsInput, setTagsInput] = useState("");
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  // -------------------------------------------------------------------
  // Data loading
  // -------------------------------------------------------------------

  const loadPosts = useCallback(async () => {
    setIsLoading(true);
    try {
      const res = await fetchBlogPosts({ page, limit: LIMIT });
      setPosts(res.data);
      setTotal(res.total);
    } catch {
      setPosts([]);
      setTotal(0);
    } finally {
      setIsLoading(false);
    }
  }, [page]);

  useEffect(() => {
    loadPosts();
  }, [loadPosts]);

  const totalPages = Math.max(1, Math.ceil(total / LIMIT));

  // -------------------------------------------------------------------
  // Form helpers
  // -------------------------------------------------------------------

  function openCreateForm() {
    setEditingId(null);
    setForm(emptyForm());
    setTagsInput("");
    setFormError(null);
    setView("form");
  }

  function openEditForm(post: BlogPost) {
    setEditingId(post.id);
    setForm({
      title: post.title,
      slug: post.slug,
      excerpt: post.excerpt,
      content: post.content,
      cover_image_url: post.cover_image_url,
      meta_title: post.meta_title,
      meta_description: post.meta_description,
      tags: post.tags,
      status: post.status,
    });
    setTagsInput(post.tags.join(", "));
    setFormError(null);
    setView("form");
  }

  function cancelForm() {
    setView("list");
    setEditingId(null);
    setForm(emptyForm());
    setTagsInput("");
    setFormError(null);
  }

  function handleTitleChange(title: string) {
    setForm((prev) => ({
      ...prev,
      title,
      slug: editingId ? prev.slug : slugify(title),
    }));
  }

  function handleTagsInputChange(value: string) {
    setTagsInput(value);
    setForm((prev) => ({
      ...prev,
      tags: value
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean),
    }));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setFormError(null);
    setSaving(true);

    try {
      if (editingId) {
        await updateBlogPost(editingId, form);
      } else {
        await createBlogPost(form);
      }
      setView("list");
      setEditingId(null);
      setForm(emptyForm());
      setTagsInput("");
      await loadPosts();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "Failed to save post");
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Are you sure you want to delete this post?")) return;

    try {
      await deleteBlogPost(id);
      await loadPosts();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to delete post");
    }
  }

  // -------------------------------------------------------------------
  // Table columns
  // -------------------------------------------------------------------

  const columns: Column<BlogPost>[] = [
    {
      key: "title",
      header: "Title",
      render: (p) => (
        <span className="font-medium line-clamp-1">{p.title}</span>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (p) => (
        <span
          className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${
            p.status === "published"
              ? "bg-success/10 text-success"
              : "bg-warning/10 text-warning"
          }`}
        >
          {p.status}
        </span>
      ),
    },
    {
      key: "tags",
      header: "Tags",
      render: (p) => (
        <div className="flex flex-wrap gap-1">
          {p.tags.map((tag) => (
            <span
              key={tag}
              className="rounded-md bg-muted px-2 py-0.5 text-xs text-muted-foreground"
            >
              {tag}
            </span>
          ))}
        </div>
      ),
    },
    {
      key: "created_at",
      header: "Created",
      render: (p) =>
        new Date(p.created_at).toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        }),
    },
    {
      key: "actions",
      header: "",
      render: (p) => (
        <div className="flex items-center gap-2 justify-end">
          <button
            onClick={() => openEditForm(p)}
            className="rounded-lg border border-border px-3 py-1 text-xs text-card-foreground transition-colors hover:bg-muted"
          >
            Edit
          </button>
          <button
            onClick={() => handleDelete(p.id)}
            className="rounded-lg border border-destructive/30 px-3 py-1 text-xs text-destructive transition-colors hover:bg-destructive/10"
          >
            Delete
          </button>
        </div>
      ),
    },
  ];

  // -------------------------------------------------------------------
  // Render
  // -------------------------------------------------------------------

  if (view === "form") {
    return (
      <div className="mx-auto max-w-3xl space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold text-foreground">
            {editingId ? "Edit Post" : "New Post"}
          </h1>
          <button
            onClick={cancelForm}
            className="rounded-lg border border-border px-4 py-2 text-sm text-muted-foreground transition-colors hover:bg-muted"
          >
            Cancel
          </button>
        </div>

        {formError && (
          <div className="rounded-lg bg-destructive/10 px-4 py-3 text-sm text-destructive">
            {formError}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          {/* Title */}
          <FormField label="Title">
            <input
              type="text"
              required
              value={form.title}
              onChange={(e) => handleTitleChange(e.target.value)}
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="Post title"
            />
          </FormField>

          {/* Slug */}
          <FormField label="Slug">
            <input
              type="text"
              required
              value={form.slug}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, slug: e.target.value }))
              }
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="post-slug"
            />
          </FormField>

          {/* Excerpt */}
          <FormField label="Excerpt">
            <textarea
              value={form.excerpt}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, excerpt: e.target.value }))
              }
              rows={2}
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring resize-y"
              placeholder="Short description..."
            />
          </FormField>

          {/* Content */}
          <FormField label="Content">
            <textarea
              value={form.content}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, content: e.target.value }))
              }
              rows={12}
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring resize-y font-mono"
              placeholder="Write your post content here..."
            />
          </FormField>

          {/* Cover Image URL */}
          <FormField label="Cover Image URL">
            <input
              type="text"
              value={form.cover_image_url}
              onChange={(e) =>
                setForm((prev) => ({
                  ...prev,
                  cover_image_url: e.target.value,
                }))
              }
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="https://example.com/image.jpg"
            />
          </FormField>

          {/* Meta title */}
          <FormField label="Meta Title">
            <input
              type="text"
              value={form.meta_title}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, meta_title: e.target.value }))
              }
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="SEO title"
            />
          </FormField>

          {/* Meta description */}
          <FormField label="Meta Description">
            <textarea
              value={form.meta_description}
              onChange={(e) =>
                setForm((prev) => ({
                  ...prev,
                  meta_description: e.target.value,
                }))
              }
              rows={2}
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring resize-y"
              placeholder="SEO description"
            />
          </FormField>

          {/* Tags */}
          <FormField label="Tags (comma-separated)">
            <input
              type="text"
              value={tagsInput}
              onChange={(e) => handleTagsInputChange(e.target.value)}
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="react, typescript, tutorial"
            />
          </FormField>

          {/* Status */}
          <FormField label="Status">
            <select
              value={form.status}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, status: e.target.value }))
              }
              className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            >
              <option value="draft">Draft</option>
              <option value="published">Published</option>
            </select>
          </FormField>

          {/* Submit */}
          <div className="flex gap-3 pt-2">
            <button
              type="submit"
              disabled={saving}
              className="rounded-lg bg-primary px-6 py-2.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving
                ? "Saving..."
                : editingId
                  ? "Update Post"
                  : "Create Post"}
            </button>
            <button
              type="button"
              onClick={cancelForm}
              className="rounded-lg border border-border px-6 py-2.5 text-sm text-muted-foreground transition-colors hover:bg-muted"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-foreground">Blog Posts</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Create and manage blog content
          </p>
        </div>

        <button
          onClick={openCreateForm}
          className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          New Post
        </button>
      </div>

      <DataTable
        columns={columns}
        data={posts}
        keyExtractor={(p) => p.id}
        page={page}
        totalPages={totalPages}
        onPageChange={setPage}
        isLoading={isLoading}
      />
    </div>
  );
}

// ---------------------------------------------------------------------------
// Form field wrapper
// ---------------------------------------------------------------------------

function FormField({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <label className="mb-1.5 block text-sm font-medium text-card-foreground">
        {label}
      </label>
      {children}
    </div>
  );
}
