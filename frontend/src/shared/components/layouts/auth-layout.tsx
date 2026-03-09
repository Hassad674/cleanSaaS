import Link from "next/link";
import { siteConfig } from "@/config/site";

export function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-50 dark:bg-zinc-950 px-4">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <Link href="/" className="text-2xl font-bold">
            {siteConfig.name}
          </Link>
        </div>
        {children}
      </div>
    </div>
  );
}
