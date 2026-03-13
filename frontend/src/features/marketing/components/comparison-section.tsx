const comparisons = [
  {
    feature: "Architecture",
    cleanSaaS: "Hexagonal (Go) + Feature-based (Next.js)",
    other: "MVC / flat structure",
  },
  {
    feature: "Modularity",
    cleanSaaS: "Delete any feature, zero errors",
    other: "Features tightly coupled",
  },
  {
    feature: "Database",
    cleanSaaS: "Pure SQL, no ORM overhead",
    other: "Prisma / ORM with migrations",
  },
  {
    feature: "AI Integration",
    cleanSaaS: "Gemini streaming + swappable adapter",
    other: "No AI or basic wrapper",
  },
  {
    feature: "Admin Panel",
    cleanSaaS: "Separate Vite app with CMS",
    other: "Shared app or none",
  },
] as const;

export function ComparisonSection() {
  return (
    <section className="py-16 sm:py-24">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
            Built different
          </h2>
          <p className="text-muted-foreground text-lg">
            CleanSaaS vs typical SaaS starters
          </p>
        </div>
        <div className="bg-card border border-border rounded-xl overflow-hidden">
          <div className="grid grid-cols-3 bg-muted px-6 py-3 text-sm font-medium text-foreground">
            <div>Feature</div>
            <div>CleanSaaS</div>
            <div>Typical Starter</div>
          </div>
          {comparisons.map((row) => (
            <div
              key={row.feature}
              className="grid grid-cols-3 px-6 py-4 border-t border-border text-sm"
            >
              <div className="font-medium text-foreground">{row.feature}</div>
              <div className="text-foreground">{row.cleanSaaS}</div>
              <div className="text-muted-foreground">{row.other}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
