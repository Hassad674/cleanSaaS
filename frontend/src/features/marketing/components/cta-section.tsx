import Link from "next/link";

export function CTASection() {
  return (
    <section className="py-16 sm:py-24">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-3xl text-center">
        <div className="bg-card border border-border rounded-2xl p-8 sm:p-12">
          <h2 className="text-2xl sm:text-3xl font-bold text-foreground mb-4">
            Stop rebuilding the same things
          </h2>
          <p className="text-muted-foreground text-lg mb-8 max-w-xl mx-auto">
            Focus on what makes your product unique. We handle auth, billing,
            storage, and the rest.
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-3">
            <Link
              href="/register"
              className="bg-primary text-primary-foreground px-6 py-3 rounded-lg text-sm font-medium hover:opacity-90 transition-opacity"
            >
              Start building
            </Link>
            <a
              href="https://github.com/Hassad674/cleanSaaS"
              target="_blank"
              rel="noopener noreferrer"
              className="border border-border bg-card text-foreground px-6 py-3 rounded-lg text-sm font-medium hover:bg-accent transition-colors"
            >
              Star on GitHub
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
