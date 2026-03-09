import type { Metadata } from "next";

export const metadata: Metadata = { title: "Blog" };

export default function BlogPage() {
  return (
    <div className="container mx-auto px-4 py-16">
      <h1 className="text-4xl font-bold mb-4">Blog</h1>
      <p className="text-zinc-600 dark:text-zinc-400">Posts will appear here.</p>
    </div>
  );
}
