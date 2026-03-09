import type { Metadata } from "next";

export const metadata: Metadata = { title: "Admin" };

export default function AdminPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Admin Dashboard</h1>
      <p className="text-zinc-500">Admin analytics and user management will be here.</p>
    </div>
  );
}
