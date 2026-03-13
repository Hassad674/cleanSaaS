export function SpotlightAdmin() {
  return (
    <section className="py-16 sm:py-24">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div className="order-2 lg:order-1 bg-card border border-border rounded-xl p-6 shadow-sm">
            <div className="space-y-4">
              <div className="flex items-center justify-between border-b border-border pb-3">
                <h3 className="font-semibold text-foreground">Admin Dashboard</h3>
                <span className="text-xs text-muted-foreground">admin.cleansaas.dev</span>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-muted rounded-lg p-3">
                  <p className="text-xs text-muted-foreground">Total Users</p>
                  <p className="text-2xl font-bold text-foreground">1,247</p>
                </div>
                <div className="bg-muted rounded-lg p-3">
                  <p className="text-xs text-muted-foreground">Blog Posts</p>
                  <p className="text-2xl font-bold text-foreground">42</p>
                </div>
              </div>
              <div className="text-xs text-muted-foreground border-t border-border pt-3">
                User management &bull; Blog CMS &bull; Role-based access
              </div>
            </div>
          </div>
          <div className="order-1 lg:order-2">
            <span className="text-xs font-medium uppercase tracking-wide text-primary mb-2 block">
              Spotlight
            </span>
            <h2 className="text-3xl font-bold text-foreground mb-4">
              Admin panel included
            </h2>
            <p className="text-muted-foreground mb-6">
              A separate Vite + React app for managing your SaaS. User management, blog CMS,
              and dashboard stats — all protected by role-based access control.
            </p>
            <ul className="space-y-3 text-sm text-muted-foreground">
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Manage users: search, change roles, view activity
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Blog CMS with draft/publish workflow and SEO fields
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Separate deployment — same design tokens, different app
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>
  );
}
