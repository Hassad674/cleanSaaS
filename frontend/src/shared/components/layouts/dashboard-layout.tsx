"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { siteConfig } from "@/config/site";
import { cn } from "@/shared/lib/utils";
import { useAuth } from "@/shared/hooks/use-auth";

export function DashboardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { user, logout } = useAuth({ required: true });

  return (
    <div className="min-h-screen flex bg-background text-foreground">
      {/* Sidebar — hidden on mobile */}
      <aside className="w-64 border-r border-sidebar-border bg-sidebar p-4 hidden lg:flex lg:flex-col shrink-0">
        <Link href="/" className="text-xl font-bold tracking-tight block mb-8">
          {siteConfig.name}
        </Link>
        <nav className="space-y-1 flex-1">
          {siteConfig.nav.dashboard.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "block px-3 py-2 rounded-lg text-sm transition-colors",
                pathname === item.href
                  ? "bg-sidebar-accent text-sidebar-accent-foreground font-medium"
                  : "text-sidebar-foreground/60 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground"
              )}
            >
              {item.label}
            </Link>
          ))}
        </nav>
        {user && (
          <div className="border-t border-sidebar-border pt-4 mt-4">
            <p className="text-sm font-medium text-sidebar-foreground px-3 truncate">
              {user.name}
            </p>
            <p className="text-xs text-sidebar-foreground/60 px-3 truncate mb-2">
              {user.email}
            </p>
            <button
              onClick={logout}
              className="w-full text-left px-3 py-2 rounded-lg text-sm text-sidebar-foreground/60 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground transition-colors"
            >
              Log out
            </button>
          </div>
        )}
      </aside>

      {/* Mobile top bar */}
      <div className="flex-1 flex flex-col min-w-0">
        <header className="lg:hidden border-b border-border bg-background px-4 h-14 flex items-center justify-between shrink-0">
          <Link href="/" className="text-lg font-bold tracking-tight">
            {siteConfig.name}
          </Link>
          {user && (
            <button
              onClick={logout}
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Log out
            </button>
          )}
        </header>
        <main className="flex-1 p-4 sm:p-6 lg:p-8 overflow-auto">{children}</main>
      </div>
    </div>
  );
}
