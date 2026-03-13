import Link from "next/link";

export function Hero() {
  return (
    <section className="py-20 sm:py-28 lg:py-36">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl text-center">
        <div className="inline-flex items-center gap-2 bg-accent text-accent-foreground px-3 py-1 rounded-full text-sm font-medium mb-6">
          <span className="w-2 h-2 rounded-full bg-primary" />
          Open Source — MIT License
        </div>
        <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-foreground mb-6 leading-tight">
          Ship your SaaS
          <span className="text-primary"> in weeks,</span>
          <br />
          not months
        </h1>
        <p className="text-lg sm:text-xl text-muted-foreground mb-10 max-w-2xl mx-auto leading-relaxed">
          A production-ready boilerplate with auth, billing, AI, storage, and
          more. Next.js + Go + PostgreSQL. Every module is independent —
          use only what you need.
        </p>
        <div className="flex flex-col sm:flex-row justify-center gap-3">
          <Link
            href="/register"
            className="bg-primary text-primary-foreground px-6 py-3 rounded-lg text-sm font-medium hover:opacity-90 transition-opacity"
          >
            Get started free
          </Link>
          <a
            href="https://github.com/Hassad674/cleanSaaS"
            target="_blank"
            rel="noopener noreferrer"
            className="border border-border bg-card text-foreground px-6 py-3 rounded-lg text-sm font-medium hover:bg-accent transition-colors"
          >
            View on GitHub
          </a>
        </div>
      </div>
    </section>
  );
}
