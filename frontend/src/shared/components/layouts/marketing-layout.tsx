import Link from "next/link";
import { siteConfig } from "@/config/site";

export function MarketingLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-zinc-200 dark:border-zinc-800">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <Link href="/" className="text-xl font-bold">
            {siteConfig.name}
          </Link>
          <nav className="hidden md:flex items-center gap-6">
            {siteConfig.nav.marketing.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="text-sm text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
              >
                {item.label}
              </Link>
            ))}
          </nav>
          <div className="flex items-center gap-3">
            <Link
              href="/login"
              className="text-sm text-zinc-600 hover:text-zinc-900 dark:text-zinc-400"
            >
              Log in
            </Link>
            <Link
              href="/register"
              className="text-sm bg-zinc-900 text-white px-4 py-2 rounded-md hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              Get started
            </Link>
          </div>
        </div>
      </header>
      <main className="flex-1">{children}</main>
      <footer className="border-t border-zinc-200 dark:border-zinc-800 py-8">
        <div className="container mx-auto px-4 text-center text-sm text-zinc-500">
          {siteConfig.name} — Open Source SaaS Boilerplate
        </div>
      </footer>
    </div>
  );
}
