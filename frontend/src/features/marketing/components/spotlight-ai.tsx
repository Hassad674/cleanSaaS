export function SpotlightAI() {
  return (
    <section className="py-16 sm:py-24 bg-muted/30">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-5xl">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <span className="text-xs font-medium uppercase tracking-wide text-primary mb-2 block">
              Spotlight
            </span>
            <h2 className="text-3xl font-bold text-foreground mb-4">
              AI Chat, ready to go
            </h2>
            <p className="text-muted-foreground mb-6">
              Ship AI-powered conversations out of the box. Gemini integration with SSE streaming,
              persistent conversation history, and auto-generated titles.
            </p>
            <ul className="space-y-3 text-sm text-muted-foreground">
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Real-time streaming responses via Server-Sent Events
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Full conversation history stored in PostgreSQL
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">&#10003;</span>
                Swap AI provider by changing one adapter file
              </li>
            </ul>
          </div>
          <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
            <div className="space-y-4">
              <div className="flex gap-3">
                <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-medium text-muted-foreground shrink-0">U</div>
                <div className="bg-muted rounded-lg p-3 text-sm text-foreground">How do I add a new feature to CleanSaaS?</div>
              </div>
              <div className="flex gap-3">
                <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-medium text-primary shrink-0">AI</div>
                <div className="bg-primary/5 border border-primary/10 rounded-lg p-3 text-sm text-foreground">
                  Create your domain entity, port interface, app service, adapter, and handler following the hexagonal architecture. Each layer only depends on the one above it...
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
