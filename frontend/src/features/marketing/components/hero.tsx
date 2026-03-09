import Link from "next/link";
import { siteConfig } from "@/config/site";

export function Hero() {
  return (
    <section className="py-24 text-center">
      <div className="container mx-auto px-4 max-w-3xl">
        <h1 className="text-5xl font-bold tracking-tight mb-6">
          The open-source SaaS boilerplate for modern apps
        </h1>
        <p className="text-xl text-zinc-600 dark:text-zinc-400 mb-8">
          {siteConfig.description}. Next.js + Go + PostgreSQL.
          Production-ready. AI-native. Deploy in minutes.
        </p>
        <div className="flex justify-center gap-4">
          <Link
            href="/register"
            className="bg-zinc-900 text-white px-6 py-3 rounded-md text-sm font-medium hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            Get started free
          </Link>
          <a
            href="https://github.com/Hassad674/cleanSaaS"
            target="_blank"
            rel="noopener noreferrer"
            className="border border-zinc-300 px-6 py-3 rounded-md text-sm font-medium hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-900"
          >
            View on GitHub
          </a>
        </div>
      </div>
    </section>
  );
}
