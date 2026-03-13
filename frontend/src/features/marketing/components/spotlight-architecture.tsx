export function SpotlightArchitecture() {
  return (
    <section className="py-16 sm:py-24 bg-muted/30">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <span className="text-xs font-medium uppercase tracking-wide text-primary mb-2 block">
              Spotlight
            </span>
            <h2 className="text-3xl font-bold text-foreground mb-4">
              Architecture that scales with you
            </h2>
            <p className="text-muted-foreground mb-6">
              Hexagonal architecture on the backend, feature-based on the frontend.
              Every module is independently removable — delete a folder, zero compilation errors.
            </p>
            <ul className="space-y-3 text-sm text-muted-foreground">
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Swap Stripe for Lemon Squeezy: one adapter file + one line in main.go
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Test business logic without a database using interface mocks
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Remove any feature without touching other code
              </li>
            </ul>
          </div>
          <div className="bg-card border border-border rounded-xl p-6 shadow-sm font-mono text-sm">
            <p className="text-muted-foreground mb-2">// Swap a provider in seconds</p>
            <div className="space-y-1 text-foreground">
              <p><span className="text-primary">-</span> paymentSvc := stripe.NewPaymentService(key)</p>
              <p><span className="text-success">+</span> paymentSvc := lemonsqueezy.NewPaymentService(key)</p>
            </div>
            <p className="text-muted-foreground mt-4 mb-2">// Remove a feature cleanly</p>
            <div className="space-y-1 text-foreground">
              <p>$ rm -rf backend/internal/app/billing/</p>
              <p>$ rm -rf frontend/src/features/billing/</p>
              <p className="text-success">$ go build ./... <span className="text-muted-foreground">// still compiles</span></p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
