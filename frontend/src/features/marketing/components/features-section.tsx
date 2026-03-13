const features = [
  {
    name: "Authentication",
    description: "Email/password, OAuth (Google, GitHub), JWT sessions, password reset. Production-ready from day one.",
    badge: "Core",
  },
  {
    name: "Billing & Subscriptions",
    description: "Stripe integration with plans, invoices, customer portal, and webhook handling.",
    badge: "Module",
  },
  {
    name: "AI Chat",
    description: "Multi-provider support — Claude, OpenAI, Gemini. Conversation history, streaming responses.",
    badge: "Module",
  },
  {
    name: "File Storage",
    description: "Upload, manage, and serve files via Cloudflare R2. Signed URLs, quotas, MIME validation.",
    badge: "Module",
  },
  {
    name: "Notifications",
    description: "Email + in-app notifications with templates. Powered by Resend. Read/unread tracking.",
    badge: "Module",
  },
  {
    name: "Admin Dashboard",
    description: "User management, analytics overview, audit logs. Role-based access control.",
    badge: "Module",
  },
] as const;

export function FeaturesSection() {
  return (
    <section id="features" className="py-16 sm:py-24">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
            Everything you need, nothing you don&apos;t
          </h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            Every feature is an independent module. Use all of them or pick only
            what your product needs.
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {features.map((feature) => (
            <div
              key={feature.name}
              className="bg-card border border-border rounded-xl p-6 flex flex-col"
            >
              <span
                className={`text-xs font-medium uppercase tracking-wide w-fit px-2 py-0.5 rounded-md mb-3 ${
                  feature.badge === "Core"
                    ? "bg-primary/10 text-primary"
                    : "bg-muted text-muted-foreground"
                }`}
              >
                {feature.badge}
              </span>
              <h3 className="text-foreground font-semibold text-lg mb-2">{feature.name}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed flex-1">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
