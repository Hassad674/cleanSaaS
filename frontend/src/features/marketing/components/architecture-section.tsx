const principles = [
  {
    title: "Fully modular",
    description:
      "Delete any feature folder — zero compilation errors. No cross-feature imports, no hidden dependencies. Ready for the future CLI that lets users pick modules.",
  },
  {
    title: "Hexagonal backend",
    description:
      "Domain → Ports → App → Adapters → Handlers. Swap Stripe for LemonSqueezy? One file, one line in main.go. Business logic stays untouched.",
  },
  {
    title: "Feature-based frontend",
    description:
      "Each feature owns its components, actions, hooks, and types. Pages are thin routing layers. Features never import each other.",
  },
  {
    title: "Pure SQL, no ORM",
    description:
      "Full control over queries, no magic. Parameterized queries, proper indexes, migration-based schema versioning with golang-migrate.",
  },
] as const;

export function ArchitectureSection() {
  return (
    <section className="py-16 sm:py-24 bg-muted/50">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
            Architecture that scales with you
          </h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            Not just a starter template — a professional-grade foundation designed for real products.
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {principles.map((principle) => (
            <div
              key={principle.title}
              className="bg-card border border-border rounded-xl p-6"
            >
              <h3 className="text-foreground font-semibold text-lg mb-2">
                {principle.title}
              </h3>
              <p className="text-sm text-muted-foreground leading-relaxed">
                {principle.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
