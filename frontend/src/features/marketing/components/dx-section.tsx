const items = [
  {
    title: "8 Claude Code skills",
    description:
      "Custom slash commands: /add-feature, /check, /review, /test, /add-migration, /add-endpoint, /add-adapter, /remove-feature. AI builds features for you.",
  },
  {
    title: "TDD-first workflow",
    description:
      "Every layer is testable. Domain tests, service tests with mocks, integration tests against real PostgreSQL. 80%+ coverage target on business logic.",
  },
  {
    title: "One command to extend",
    description:
      "Add a feature: domain → ports → service → adapter → handler → frontend. Consistent patterns across every module — learn one, know all.",
  },
] as const;

export function DXSection() {
  return (
    <section className="py-16 sm:py-24">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
            Built for developers who ship fast
          </h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            AI-native developer experience. Consistent patterns. Zero guesswork.
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {items.map((item) => (
            <div
              key={item.title}
              className="bg-card border border-border rounded-xl p-6"
            >
              <h3 className="text-foreground font-semibold text-lg mb-2">
                {item.title}
              </h3>
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
