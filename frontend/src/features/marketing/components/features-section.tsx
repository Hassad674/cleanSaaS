const features = [
  {
    name: "Authentication",
    description: "Email/password login, JWT sessions, password reset, email verification. Production-ready auth from day one.",
    badge: "Core",
  },
  {
    name: "Billing & Subscriptions",
    description: "Stripe Checkout, 3 plans, webhooks, invoices, customer portal. Plug in your Stripe key and go.",
    badge: "Module",
  },
  {
    name: "AI Chat",
    description: "Gemini-powered conversations with SSE streaming, conversation history, and auto-titling.",
    badge: "Module",
  },
  {
    name: "File Storage",
    description: "Upload and manage files via Cloudflare R2. Drag-and-drop UI, type validation, 50MB limit.",
    badge: "Module",
  },
  {
    name: "Email System",
    description: "Transactional emails via Resend with HTML templates for verification, password reset, and more.",
    badge: "Module",
  },
  {
    name: "Notifications",
    description: "In-app notification bell with unread count, mark-as-read, and 30s polling. Always available.",
    badge: "Module",
  },
  {
    name: "Blog CMS",
    description: "Database-backed blog with tags, slugs, SEO metadata, and draft/publish workflow. Admin-managed.",
    badge: "Module",
  },
  {
    name: "Admin Panel",
    description: "Separate Vite app for user management, blog editing, and dashboard stats. Role-based access.",
    badge: "Module",
  },
  {
    name: "Background Jobs",
    description: "Go-native job scheduler with token cleanup and system stats logging. No external dependencies.",
    badge: "Core",
  },
  {
    name: "Security",
    description: "Rate limiting, bcrypt passwords, parameterized SQL, CORS, input validation at every boundary.",
    badge: "Core",
  },
  {
    name: "Hexagonal Architecture",
    description: "Clean separation of concerns. Swap any provider by changing one adapter. Test logic without infrastructure.",
    badge: "Core",
  },
  {
    name: "Developer Experience",
    description: "CLAUDE.md at every level, skills system, 74+ tests, TypeScript strict mode, conventional commits.",
    badge: "Core",
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
            what your product needs. Delete a folder, zero compilation errors.
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
