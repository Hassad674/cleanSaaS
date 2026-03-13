import { useAuth } from "../hooks/useAuth.ts";

export default function Header() {
  const { user, logout } = useAuth();

  return (
    <header className="sticky top-0 z-20 flex h-16 items-center justify-between border-b border-border bg-card px-6">
      {/* Left: page context (can be extended) */}
      <div />

      {/* Right: user info + logout */}
      <div className="flex items-center gap-4">
        <div className="text-right">
          <p className="text-sm font-medium text-card-foreground">
            {user?.name ?? user?.email}
          </p>
          <p className="text-xs text-muted-foreground capitalize">
            {user?.role}
          </p>
        </div>

        <button
          onClick={logout}
          className="rounded-lg border border-border px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          Logout
        </button>
      </div>
    </header>
  );
}
