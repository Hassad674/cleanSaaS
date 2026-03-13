import { useEffect, useState } from "react";
import StatsCard from "../components/StatsCard.tsx";
import { fetchUsers, fetchBlogPosts } from "../lib/api.ts";

export default function DashboardPage() {
  const [totalUsers, setTotalUsers] = useState<number | null>(null);
  const [totalPosts, setTotalPosts] = useState<number | null>(null);
  const [publishedPosts, setPublishedPosts] = useState<number | null>(null);
  const [draftPosts, setDraftPosts] = useState<number | null>(null);

  useEffect(() => {
    // Fetch counts — we only need totals, so limit=1 is enough
    fetchUsers({ page: 1, limit: 1 })
      .then((res) => setTotalUsers(res.total))
      .catch(() => setTotalUsers(0));

    fetchBlogPosts({ page: 1, limit: 1 })
      .then((res) => setTotalPosts(res.total))
      .catch(() => setTotalPosts(0));

    fetchBlogPosts({ page: 1, limit: 1, status: "published" })
      .then((res) => setPublishedPosts(res.total))
      .catch(() => setPublishedPosts(0));

    fetchBlogPosts({ page: 1, limit: 1, status: "draft" })
      .then((res) => setDraftPosts(res.total))
      .catch(() => setDraftPosts(0));
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-foreground">Dashboard</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Overview of your SaaS application
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          icon={<UsersStatIcon />}
          label="Total Users"
          value={totalUsers ?? "..."}
        />
        <StatsCard
          icon={<PostsStatIcon />}
          label="Total Posts"
          value={totalPosts ?? "..."}
        />
        <StatsCard
          icon={<PublishedStatIcon />}
          label="Published"
          value={publishedPosts ?? "..."}
        />
        <StatsCard
          icon={<DraftStatIcon />}
          label="Drafts"
          value={draftPosts ?? "..."}
        />
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Icons
// ---------------------------------------------------------------------------

function UsersStatIcon() {
  return (
    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" />
    </svg>
  );
}

function PostsStatIcon() {
  return (
    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 7.5h1.5m-1.5 3h1.5m-7.5 3h7.5m-7.5 3h7.5m3-9h3.375c.621 0 1.125.504 1.125 1.125V18a2.25 2.25 0 0 1-2.25 2.25M16.5 7.5V4.875c0-.621-.504-1.125-1.125-1.125H4.125C3.504 3.75 3 4.254 3 4.875V18a2.25 2.25 0 0 0 2.25 2.25h13.5" />
    </svg>
  );
}

function PublishedStatIcon() {
  return (
    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
    </svg>
  );
}

function DraftStatIcon() {
  return (
    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
    </svg>
  );
}
