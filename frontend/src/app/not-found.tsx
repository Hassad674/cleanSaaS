import Link from "next/link";

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-6xl font-bold mb-4">404</h1>
        <p className="text-zinc-500 mb-6">Page not found</p>
        <Link
          href="/"
          className="bg-zinc-900 text-white px-4 py-2 rounded-md hover:bg-zinc-800 dark:bg-white dark:text-zinc-900"
        >
          Go home
        </Link>
      </div>
    </div>
  );
}
