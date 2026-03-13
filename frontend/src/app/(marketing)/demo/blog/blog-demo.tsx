"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface BlogPost {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  content: string[];
  tag: string;
  readTime: number;
  author: { name: string; initials: string };
  date: string;
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

const TAGS = ["All", "Guide", "Architecture", "AI", "Billing", "DevOps", "Frontend"] as const;

const TAG_COLORS: Record<string, string> = {
  Guide: "bg-primary/10 text-primary",
  Architecture: "bg-accent text-foreground",
  AI: "bg-primary/15 text-primary",
  Billing: "bg-muted text-foreground",
  DevOps: "bg-primary/10 text-primary",
  Frontend: "bg-accent text-foreground",
};

const POSTS: BlogPost[] = [
  {
    id: "1",
    title: "Getting Started with CleanSaaS",
    slug: "getting-started-with-cleansaas",
    excerpt:
      "A step-by-step guide to setting up your development environment and shipping your first feature with CleanSaaS.",
    content: [
      "CleanSaaS is designed to get you from zero to production in record time. The boilerplate ships with a complete stack — Next.js on the frontend, Go with Chi on the backend, and PostgreSQL for persistence — all wired together with sensible defaults and professional-grade patterns.",
      "Setting up locally is straightforward. Clone the repository, run docker compose up to spin up PostgreSQL and DbGate, apply migrations with make migrate-up, and start both the backend and frontend dev servers. Within minutes you have a fully functional SaaS application running on your machine with authentication, billing, file storage, and more.",
      "The real power of CleanSaaS lies in its modularity. Every feature is independently removable. If you do not need AI chat, delete its folder and remove its wiring from main.go — nothing else breaks. This architecture means you start with everything and strip away what you do not need, rather than building from scratch.",
    ],
    tag: "Guide",
    readTime: 5,
    author: { name: "John Doe", initials: "JD" },
    date: "2026-03-01",
  },
  {
    id: "2",
    title: "Why Hexagonal Architecture Matters",
    slug: "why-hexagonal-architecture-matters",
    excerpt:
      "How ports and adapters keep your Go backend testable, swappable, and ready for any infrastructure change.",
    content: [
      "Hexagonal architecture — also known as ports and adapters — is the backbone of the CleanSaaS backend. The core idea is simple: your business logic should never depend on infrastructure. Instead, the domain layer defines interfaces (ports), and concrete implementations (adapters) are injected at startup.",
      "In practice this means the billing service does not know whether payments go through Stripe or LemonSqueezy. It calls a PaymentProvider interface. Swapping providers is a one-file change: write a new adapter that implements the interface and wire it in main.go. No business logic touched, no tests rewritten.",
      "This pattern pays dividends beyond swappability. Testing becomes trivial — inject a mock adapter and test pure business logic without databases, HTTP calls, or third-party APIs. It also enforces a clean dependency graph: adapters depend on the domain, never the other way around. The result is a codebase that scales with your team and your product.",
    ],
    tag: "Architecture",
    readTime: 8,
    author: { name: "John Doe", initials: "JD" },
    date: "2026-02-22",
  },
  {
    id: "3",
    title: "Building Swappable AI Providers",
    slug: "building-swappable-ai-providers",
    excerpt:
      "Implement OpenAI, Anthropic, or any LLM behind a single interface with streaming support and graceful fallbacks.",
    content: [
      "The AI module in CleanSaaS is built around a single AIProvider interface with methods for chat completion, streaming, and model listing. The default adapter targets OpenAI, but the architecture makes it trivial to add Anthropic, Google Gemini, Mistral, or a self-hosted model.",
      "Streaming is where things get interesting. The interface returns a channel of token chunks, and the HTTP handler uses Server-Sent Events to push them to the frontend in real time. The frontend reads the stream with a custom hook that buffers tokens and renders incrementally. Because the interface is provider-agnostic, switching from GPT-4 to Claude requires zero frontend changes.",
      "Fallback strategies are built in at the application layer. If the primary provider returns an error or times out, the service can automatically retry with a secondary provider. This resilience logic lives above the adapter layer, so it works identically regardless of which providers you configure. Users never see a blank screen — they get a response from whichever model is available.",
    ],
    tag: "AI",
    readTime: 6,
    author: { name: "John Doe", initials: "JD" },
    date: "2026-02-15",
  },
  {
    id: "4",
    title: "Stripe Integration Best Practices",
    slug: "stripe-integration-best-practices",
    excerpt:
      "Webhooks, idempotency, subscription lifecycle — everything you need for production-grade billing in Go.",
    content: [
      "Billing is one of the most critical modules in any SaaS. CleanSaaS ships with a complete Stripe integration covering plan management, checkout sessions, subscription lifecycle events, and invoice handling. The implementation follows Stripe's own best practices and handles the edge cases that most tutorials skip.",
      "Webhook handling is where most integrations fall apart. CleanSaaS verifies every webhook signature, processes events idempotently using event IDs, and handles out-of-order delivery gracefully. If a subscription.updated event arrives before subscription.created (which happens more often than you think), the handler creates the record first and then applies the update. This resilience is critical for production.",
      "The billing adapter is completely isolated from the rest of the application. It implements a PaymentProvider interface and communicates with the billing service through well-defined ports. Want to switch to LemonSqueezy or Paddle? Write a new adapter, swap one line in main.go, and your subscription logic, pricing pages, and customer portal all keep working without modification.",
    ],
    tag: "Billing",
    readTime: 7,
    author: { name: "John Doe", initials: "JD" },
    date: "2026-02-08",
  },
  {
    id: "5",
    title: "Deploying to Production",
    slug: "deploying-to-production",
    excerpt:
      "Ship your CleanSaaS app to Vercel, Railway, and Neon with zero-downtime deployments and production-grade monitoring.",
    content: [
      "CleanSaaS is designed for a modern deployment stack: Vercel hosts the Next.js frontend, Railway runs the Go backend, and Neon provides serverless PostgreSQL. This combination gives you global edge delivery, automatic scaling, and branching databases — all on generous free tiers for getting started.",
      "The deployment pipeline is straightforward. Push to main and Vercel auto-deploys the frontend. The backend deploys to Railway via a Dockerfile — Railway detects the build, creates the container, and routes traffic with zero downtime. Database migrations run as a pre-deploy step using the production DATABASE_URL. Neon's branching feature lets you test migrations against a copy of production data before applying them to the real database.",
      "Monitoring and observability come next. The backend ships with structured logging via zerolog, request tracing middleware, and health check endpoints that Railway uses for readiness probes. For production monitoring, integrate with your preferred stack — Sentry for error tracking, PostHog for analytics, or Grafana Cloud for metrics. The admin panel provides a real-time view of user activity, subscription status, and system health directly from your database.",
    ],
    tag: "DevOps",
    readTime: 10,
    author: { name: "John Doe", initials: "JD" },
    date: "2026-01-30",
  },
  {
    id: "6",
    title: "Feature-Based Frontend Architecture",
    slug: "feature-based-frontend-architecture",
    excerpt:
      "Organize your Next.js app by feature, not by type. Keep components, hooks, and actions colocated for maximum developer velocity.",
    content: [
      "Most Next.js projects organize code by type: a components folder, a hooks folder, a utils folder. This works for small apps but collapses under its own weight as the codebase grows. CleanSaaS takes a different approach: organize by feature. Each feature — auth, billing, AI, blog — is a self-contained module with its own components, hooks, actions, and types.",
      "The key rule is that features never import from each other. If the billing page needs user data, the page in the app directory imports from both features and composes them. The features themselves remain independent. This means you can delete an entire feature folder and the rest of the application compiles and runs without errors.",
      "Server Actions are scoped to their feature as well. The auth feature has its own actions that call auth API endpoints. The billing feature has its own actions for billing endpoints. No cross-feature API calls happen inside a feature. This isolation makes each feature easy to understand, test, and maintain — and it prepares the codebase for the eventual CLI tool that lets users pick only the modules they need.",
    ],
    tag: "Frontend",
    readTime: 5,
    author: { name: "John Doe", initials: "JD" },
    date: "2026-01-22",
  },
];

// ---------------------------------------------------------------------------
// Helper: format date
// ---------------------------------------------------------------------------

function formatDate(dateStr: string): string {
  return new Intl.DateTimeFormat("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
  }).format(new Date(dateStr));
}

// ---------------------------------------------------------------------------
// Component
// ---------------------------------------------------------------------------

export function BlogDemo() {
  const [activeTag, setActiveTag] = useState<string>("All");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedPost, setSelectedPost] = useState<BlogPost | null>(null);
  const [copyFeedback, setCopyFeedback] = useState(false);

  const filteredPosts = useMemo(() => {
    return POSTS.filter((post) => {
      const matchesTag = activeTag === "All" || post.tag === activeTag;
      const matchesSearch =
        searchQuery.trim() === "" ||
        post.title.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesTag && matchesSearch;
    });
  }, [activeTag, searchQuery]);

  function handleCopyLink() {
    setCopyFeedback(true);
    setTimeout(() => setCopyFeedback(false), 2000);
  }

  // -------------------------------------------------------------------------
  // Full post view
  // -------------------------------------------------------------------------

  if (selectedPost) {
    return (
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Navigation */}
        <div className="flex items-center gap-4 mb-8">
          <button
            onClick={() => setSelectedPost(null)}
            className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <svg
              className="h-4 w-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
              />
            </svg>
            Back to posts
          </button>
        </div>

        {/* Article */}
        <article className="max-w-2xl mx-auto">
          {/* Tag */}
          <span
            className={cn(
              "inline-block text-xs font-medium px-2.5 py-0.5 rounded-full mb-4",
              TAG_COLORS[selectedPost.tag] ?? "bg-muted text-foreground"
            )}
          >
            {selectedPost.tag}
          </span>

          {/* Title */}
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight leading-tight">
            {selectedPost.title}
          </h1>

          {/* Meta */}
          <div className="flex items-center gap-4 mt-6 pb-8 border-b border-border">
            <div className="h-10 w-10 rounded-full bg-primary/10 text-primary flex items-center justify-center text-sm font-semibold shrink-0">
              {selectedPost.author.initials}
            </div>
            <div className="flex flex-col sm:flex-row sm:items-center sm:gap-3">
              <span className="text-sm font-medium text-foreground">
                {selectedPost.author.name}
              </span>
              <span className="hidden sm:inline text-muted-foreground">
                &middot;
              </span>
              <span className="text-sm text-muted-foreground">
                {formatDate(selectedPost.date)}
              </span>
              <span className="hidden sm:inline text-muted-foreground">
                &middot;
              </span>
              <span className="text-sm text-muted-foreground">
                {selectedPost.readTime} min read
              </span>
            </div>
          </div>

          {/* Content */}
          <div className="mt-8 space-y-6">
            {selectedPost.content.map((paragraph, i) => (
              <p
                key={i}
                className="text-base leading-relaxed text-foreground/90"
              >
                {paragraph}
              </p>
            ))}
          </div>

          {/* Share buttons */}
          <div className="mt-12 pt-8 border-t border-border">
            <p className="text-sm font-medium text-foreground mb-4">
              Share this article
            </p>
            <div className="flex items-center gap-3">
              {/* Twitter / X */}
              <button
                onClick={() => {}}
                className="inline-flex items-center gap-2 text-sm text-muted-foreground bg-muted hover:bg-accent hover:text-foreground rounded-lg px-4 py-2 transition-colors"
              >
                <svg
                  className="h-4 w-4"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
                </svg>
                Twitter
              </button>

              {/* LinkedIn */}
              <button
                onClick={() => {}}
                className="inline-flex items-center gap-2 text-sm text-muted-foreground bg-muted hover:bg-accent hover:text-foreground rounded-lg px-4 py-2 transition-colors"
              >
                <svg
                  className="h-4 w-4"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 01-2.063-2.065 2.064 2.064 0 112.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
                </svg>
                LinkedIn
              </button>

              {/* Copy link */}
              <button
                onClick={handleCopyLink}
                className={cn(
                  "inline-flex items-center gap-2 text-sm rounded-lg px-4 py-2 transition-colors",
                  copyFeedback
                    ? "bg-primary/10 text-primary"
                    : "text-muted-foreground bg-muted hover:bg-accent hover:text-foreground"
                )}
              >
                {copyFeedback ? (
                  <>
                    <svg
                      className="h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      strokeWidth={2}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M4.5 12.75l6 6 9-13.5"
                      />
                    </svg>
                    Copied!
                  </>
                ) : (
                  <>
                    <svg
                      className="h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      strokeWidth={2}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m9.86-2.54a4.5 4.5 0 00-1.242-7.244l-4.5-4.5a4.5 4.5 0 00-6.364 6.364L4.25 8.81"
                      />
                    </svg>
                    Copy link
                  </>
                )}
              </button>
            </div>
          </div>
        </article>

        {/* Demo notice */}
        <p className="text-center text-sm text-muted-foreground mt-16">
          This is a simulated blog post. In production, content is managed from
          the admin panel and served from the database.
        </p>
      </div>
    );
  }

  // -------------------------------------------------------------------------
  // List view
  // -------------------------------------------------------------------------

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Back to demos */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <svg
          className="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
          />
        </svg>
        Back to demos
      </Link>

      {/* Header */}
      <div className="max-w-2xl mb-12">
        <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
          Blog
        </h1>
        <p className="text-lg text-muted-foreground mt-3">
          Insights on architecture, deployment, and building scalable SaaS
          products with CleanSaaS.
        </p>
      </div>

      {/* Search bar */}
      <div className="relative max-w-md mb-8">
        <svg
          className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
          />
        </svg>
        <input
          type="text"
          placeholder="Search posts..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 text-sm bg-card border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-shadow"
        />
      </div>

      {/* Tag filter */}
      <div className="flex flex-wrap items-center gap-2 mb-10">
        {TAGS.map((tag) => (
          <button
            key={tag}
            onClick={() => setActiveTag(tag)}
            className={cn(
              "text-sm font-medium px-3.5 py-1.5 rounded-full transition-colors",
              activeTag === tag
                ? "bg-primary text-primary-foreground"
                : "bg-muted text-muted-foreground hover:bg-accent hover:text-foreground"
            )}
          >
            {tag}
          </button>
        ))}
      </div>

      {/* Post grid */}
      {filteredPosts.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-muted-foreground">
            No posts found matching your criteria.
          </p>
          <button
            onClick={() => {
              setActiveTag("All");
              setSearchQuery("");
            }}
            className="text-sm text-primary hover:underline mt-2"
          >
            Clear filters
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {filteredPosts.map((post) => (
            <button
              key={post.id}
              onClick={() => setSelectedPost(post)}
              className="group text-left bg-card border border-border rounded-xl p-6 shadow-sm hover:shadow-md hover:border-primary/50 transition-all duration-200"
            >
              {/* Tag badge */}
              <span
                className={cn(
                  "inline-block text-xs font-medium px-2.5 py-0.5 rounded-full mb-3",
                  TAG_COLORS[post.tag] ?? "bg-muted text-foreground"
                )}
              >
                {post.tag}
              </span>

              {/* Title */}
              <h2 className="text-lg font-semibold text-foreground group-hover:text-primary transition-colors leading-snug">
                {post.title}
              </h2>

              {/* Excerpt */}
              <p className="text-sm text-muted-foreground mt-2 line-clamp-2 leading-relaxed">
                {post.excerpt}
              </p>

              {/* Footer: author, date, read time */}
              <div className="flex items-center gap-3 mt-5 pt-4 border-t border-border">
                {/* Avatar */}
                <div className="h-8 w-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold shrink-0">
                  {post.author.initials}
                </div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="font-medium text-foreground">
                    {post.author.name}
                  </span>
                  <span>&middot;</span>
                  <span>{formatDate(post.date)}</span>
                  <span>&middot;</span>
                  <span>{post.readTime} min read</span>
                </div>
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Demo notice */}
      <p className="text-center text-sm text-muted-foreground mt-12">
        All posts are simulated. In production, blog content is managed from the
        admin panel and stored in PostgreSQL.
      </p>
    </div>
  );
}
