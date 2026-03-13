const stack = [
  {
    name: "Next.js",
    description: "React framework with Server Components",
    category: "Frontend",
  },
  {
    name: "Go + Chi",
    description: "Fast, compiled backend with minimal dependencies",
    category: "Backend",
  },
  {
    name: "PostgreSQL",
    description: "Pure SQL, no ORM — full control over queries",
    category: "Database",
  },
  {
    name: "Tailwind CSS",
    description: "Utility-first styling with design tokens",
    category: "Styling",
  },
  {
    name: "Vercel",
    description: "Zero-config frontend deployment",
    category: "Deploy",
  },
  {
    name: "Railway",
    description: "One-click Go backend hosting",
    category: "Deploy",
  },
  {
    name: "Neon",
    description: "Serverless PostgreSQL, scales to zero",
    category: "Database",
  },
  {
    name: "Cloudflare R2",
    description: "S3-compatible storage, no egress fees",
    category: "Storage",
  },
] as const;

export function StackSection() {
  return (
    <section className="py-16 sm:py-24 bg-muted/50">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
            Modern stack, no compromise
          </h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            Every technology chosen for performance, developer experience, and production readiness.
          </p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {stack.map((item) => (
            <div
              key={item.name}
              className="bg-card border border-border rounded-xl p-5 hover:shadow-sm transition-shadow"
            >
              <span className="text-xs font-medium text-primary uppercase tracking-wide">
                {item.category}
              </span>
              <h3 className="text-foreground font-semibold mt-1 mb-1">{item.name}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">
                {item.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
