import Link from "next/link";
import { siteConfig } from "@/config/site";

export function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-muted px-4 py-12">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <Link href="/" className="text-2xl font-bold tracking-tight">
            {siteConfig.name}
          </Link>
        </div>
        <div className="bg-card border border-border rounded-xl p-6 sm:p-8 shadow-sm">
          {children}
        </div>
      </div>
    </div>
  );
}
