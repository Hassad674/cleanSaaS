"use client";

import { useState, useEffect, useCallback } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { cn, formatCurrency } from "@/shared/lib/utils";
import { api } from "@/shared/lib/api";
import type { Plan } from "@/features/billing/types";

/* -------------------------------------------------------------------------- */
/*  Types                                                                      */
/* -------------------------------------------------------------------------- */

type BillingInterval = "month" | "year" | "lifetime";

interface DemoSession {
  plan_name: string;
  price_cents: number;
  interval: string;
  customer_id: string;
  customer_email: string;
  mode: string;
  status: string;
}

/* -------------------------------------------------------------------------- */
/*  Icons                                                                      */
/* -------------------------------------------------------------------------- */

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M5 13l4 4L19 7"
      />
    </svg>
  );
}

/* -------------------------------------------------------------------------- */
/*  Interval selector                                                          */
/* -------------------------------------------------------------------------- */

function IntervalSelector({
  value,
  onChange,
}: {
  value: BillingInterval;
  onChange: (interval: BillingInterval) => void;
}) {
  const intervals: { key: BillingInterval; label: string; badge?: string }[] = [
    { key: "month", label: "Monthly" },
    { key: "year", label: "Yearly", badge: "Save ~17%" },
    { key: "lifetime", label: "Lifetime", badge: "Pay once" },
  ];

  return (
    <div className="flex items-center justify-center">
      <div className="inline-flex bg-muted rounded-lg p-1 gap-1">
        {intervals.map((interval) => (
          <button
            key={interval.key}
            type="button"
            onClick={() => onChange(interval.key)}
            className={cn(
              "relative rounded-md px-4 py-2 text-sm font-medium transition-all",
              value === interval.key
                ? "bg-card text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            {interval.label}
            {interval.badge && value === interval.key && (
              <span className="ml-1.5 text-xs font-medium text-primary bg-primary/10 px-1.5 py-0.5 rounded-full">
                {interval.badge}
              </span>
            )}
          </button>
        ))}
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Plan grouping                                                              */
/* -------------------------------------------------------------------------- */

type PlanTier = {
  tierName: string;
  plans: Record<BillingInterval, Plan | undefined>;
  features: string[];
  highlighted: boolean;
};

function groupPlansByTier(plans: Plan[]): PlanTier[] {
  const tiers: Record<string, PlanTier> = {};

  for (const plan of plans) {
    let tierName = plan.name;
    for (const suffix of [" Monthly", " Yearly", " Lifetime"]) {
      if (tierName.endsWith(suffix)) {
        tierName = tierName.slice(0, -suffix.length);
        break;
      }
    }

    if (!tiers[tierName]) {
      tiers[tierName] = {
        tierName,
        plans: { month: undefined, year: undefined, lifetime: undefined },
        features: plan.features,
        highlighted: tierName === "Pro",
      };
    }

    tiers[tierName].plans[plan.interval] = plan;
    if (plan.features.length > tiers[tierName].features.length) {
      tiers[tierName].features = plan.features;
    }
  }

  return Object.values(tiers).sort((a, b) => {
    const order = ["Free", "Pro", "Enterprise"];
    return order.indexOf(a.tierName) - order.indexOf(b.tierName);
  });
}

/* -------------------------------------------------------------------------- */
/*  Pricing card                                                               */
/* -------------------------------------------------------------------------- */

function PricingCard({
  tier,
  interval,
  onCheckout,
  loading,
  currentPlanName,
}: {
  tier: PlanTier;
  interval: BillingInterval;
  onCheckout: (planID: string) => void;
  loading: string | null;
  currentPlanName?: string;
}) {
  const plan = tier.plans[interval] ?? tier.plans["month"];
  const isFree = tier.tierName === "Free";
  const priceCents = plan?.price_cents ?? 0;
  const isLifetime = interval === "lifetime";
  const isLoading = loading === plan?.id;

  // Check if this tier matches the current plan
  const isCurrent = currentPlanName
    ? currentPlanName.startsWith(tier.tierName) && !isFree
    : false;

  const displayFeatures = isFree
    ? (tier.plans["month"]?.features ?? tier.features)
    : tier.features;

  return (
    <div
      className={cn(
        "bg-card border rounded-xl p-6 shadow-sm flex flex-col relative",
        isCurrent
          ? "border-success ring-2 ring-success"
          : tier.highlighted
            ? "border-primary ring-2 ring-primary"
            : "border-border"
      )}
    >
      {isCurrent && (
        <span className="absolute -top-3 left-1/2 -translate-x-1/2 bg-success text-primary-foreground text-xs font-medium px-3 py-1 rounded-full">
          Current plan
        </span>
      )}
      {!isCurrent && tier.highlighted && (
        <span className="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary text-primary-foreground text-xs font-medium px-3 py-1 rounded-full">
          Most popular
        </span>
      )}

      <h3 className="text-lg font-semibold text-foreground">
        {tier.tierName}
      </h3>

      <div className="mt-4 mb-6">
        <span className="text-4xl font-bold text-foreground">
          {isFree ? "Free" : formatCurrency(priceCents, "usd")}
        </span>
        {!isFree && !isLifetime && (
          <span className="text-muted-foreground ml-1">
            /{interval === "month" ? "mo" : "yr"}
          </span>
        )}
        {!isFree && isLifetime && (
          <span className="text-muted-foreground ml-1">one-time</span>
        )}
        {!isFree && interval === "year" && (
          <span className="block text-xs text-muted-foreground mt-1">
            {formatCurrency(Math.round(priceCents / 12), "usd")}/mo equivalent
          </span>
        )}
      </div>

      <ul className="space-y-3 mb-8 flex-1">
        {displayFeatures.map((feature) => (
          <li
            key={feature}
            className="flex items-start gap-2 text-sm text-foreground"
          >
            <CheckIcon className="h-5 w-5 shrink-0 text-primary mt-0.5" />
            {feature}
          </li>
        ))}
      </ul>

      {isFree ? (
        <button
          disabled
          className="w-full bg-muted text-muted-foreground rounded-lg px-4 py-2.5 font-medium cursor-default"
        >
          Free forever
        </button>
      ) : isCurrent ? (
        <button
          disabled
          className="w-full bg-success/10 text-success rounded-lg px-4 py-2.5 font-medium cursor-default"
        >
          Active
        </button>
      ) : currentPlanName ? (
        <button
          onClick={() => plan && onCheckout(plan.id)}
          disabled={isLoading || !plan}
          className="w-full bg-primary text-primary-foreground rounded-lg px-4 py-2.5 font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {isLoading ? "Redirecting..." : "Switch to this plan"}
        </button>
      ) : (
        <button
          onClick={() => plan && onCheckout(plan.id)}
          disabled={isLoading || !plan}
          className="w-full bg-primary text-primary-foreground rounded-lg px-4 py-2.5 font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {isLoading ? "Redirecting to Stripe..." : "Subscribe with Stripe"}
        </button>
      )}
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Active subscription card                                                   */
/* -------------------------------------------------------------------------- */

function ActiveSubscriptionCard({
  session,
  onManage,
  managingPortal,
}: {
  session: DemoSession;
  onManage: () => void;
  managingPortal: boolean;
}) {
  const intervalLabel =
    session.interval === "month"
      ? "Monthly"
      : session.interval === "year"
        ? "Yearly"
        : "Lifetime";

  const isLifetime = session.interval === "lifetime";

  return (
    <div className="bg-card border border-success/30 rounded-xl p-6 shadow-sm mb-8 max-w-2xl mx-auto">
      <div className="flex items-start gap-4">
        <div className="h-12 w-12 rounded-full bg-success/10 flex items-center justify-center shrink-0">
          <CheckIcon className="h-7 w-7 text-success" />
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold text-foreground">
            You&apos;re subscribed!
          </h3>
          <p className="text-sm text-muted-foreground mt-1">
            {session.customer_email && (
              <span className="text-foreground">{session.customer_email}</span>
            )}
          </p>

          <div className="mt-4 flex flex-wrap gap-3">
            <span className="inline-flex items-center gap-1.5 bg-primary/10 text-primary text-sm font-medium px-3 py-1.5 rounded-lg">
              {session.plan_name}
            </span>
            <span className="inline-flex items-center gap-1.5 bg-muted text-foreground text-sm font-medium px-3 py-1.5 rounded-lg">
              {formatCurrency(session.price_cents, "usd")}
              {!isLifetime && (
                <span className="text-muted-foreground">
                  /{session.interval === "month" ? "mo" : "yr"}
                </span>
              )}
              {isLifetime && (
                <span className="text-muted-foreground">one-time</span>
              )}
            </span>
            <span className="inline-flex items-center gap-1.5 bg-muted text-foreground text-sm font-medium px-3 py-1.5 rounded-lg">
              {intervalLabel}
            </span>
          </div>

          <div className="mt-5 flex flex-wrap gap-3">
            <button
              onClick={onManage}
              disabled={managingPortal}
              className="inline-flex items-center gap-2 bg-primary text-primary-foreground text-sm font-medium px-4 py-2 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              <svg
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={2}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75"
                />
              </svg>
              {managingPortal
                ? "Opening portal..."
                : "Manage subscription"}
            </button>
            <p className="text-xs text-muted-foreground self-center">
              Upgrade, downgrade, cancel, or update payment method via Stripe
              portal
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Cancel banner                                                              */
/* -------------------------------------------------------------------------- */

function CancelBanner() {
  return (
    <div className="bg-muted border border-border rounded-xl p-6 mb-8 max-w-2xl mx-auto">
      <div className="flex items-start gap-3">
        <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center shrink-0">
          <svg
            className="h-6 w-6 text-muted-foreground"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </div>
        <div>
          <h3 className="text-lg font-semibold text-foreground">
            Checkout canceled
          </h3>
          <p className="text-sm text-muted-foreground mt-1">
            You canceled the Stripe checkout. Feel free to try again with any
            plan.
          </p>
        </div>
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Main demo component                                                        */
/* -------------------------------------------------------------------------- */

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

export function BillingDemo() {
  const searchParams = useSearchParams();
  const success = searchParams.get("success") === "true";
  const canceled = searchParams.get("canceled") === "true";
  const sessionId = searchParams.get("session_id");

  const [interval, setInterval] = useState<BillingInterval>("month");
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [checkoutLoading, setCheckoutLoading] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [demoSession, setDemoSession] = useState<DemoSession | null>(null);
  const [managingPortal, setManagingPortal] = useState(false);

  // Fetch plans
  useEffect(() => {
    api<Plan[]>("/billing/plans").then((res) => {
      if (res.data) {
        setPlans(res.data);
      } else {
        setError(res.error ?? "Failed to load plans");
      }
      setLoading(false);
    });
  }, []);

  // Fetch session details after successful checkout
  useEffect(() => {
    if (success && sessionId) {
      fetch(`${API_URL}/demo/billing/session?session_id=${sessionId}`)
        .then((res) => res.json())
        .then((data: DemoSession) => {
          if (data.plan_name) {
            setDemoSession(data);
            // Set interval selector to match the active plan
            if (
              data.interval === "month" ||
              data.interval === "year" ||
              data.interval === "lifetime"
            ) {
              setInterval(data.interval);
            }
          }
        })
        .catch(() => {
          // Silently fail — still show success state
        });
    }
  }, [success, sessionId]);

  const tiers = groupPlansByTier(plans);

  const handleCheckout = useCallback(async (planID: string) => {
    setCheckoutLoading(planID);
    setError(null);

    const currentURL = window.location.origin + window.location.pathname;
    const res = await api<{ url: string }>("/demo/billing/checkout", {
      method: "POST",
      body: {
        plan_id: planID,
        success_url: `${currentURL}?success=true&session_id={CHECKOUT_SESSION_ID}`,
        cancel_url: `${currentURL}?canceled=true`,
      },
    });

    if (res.data?.url) {
      window.location.href = res.data.url;
    } else {
      setError(res.error ?? "Failed to create checkout session");
      setCheckoutLoading(null);
    }
  }, []);

  const handleManagePortal = useCallback(async () => {
    if (!demoSession?.customer_id) return;
    setManagingPortal(true);

    const returnURL = window.location.href;
    try {
      const res = await fetch(`${API_URL}/demo/billing/portal`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          customer_id: demoSession.customer_id,
          return_url: returnURL,
        }),
      });
      const data = await res.json();
      if (data.url) {
        window.location.href = data.url;
      } else {
        setError("Failed to open billing portal");
      }
    } catch {
      setError("Failed to open billing portal");
    } finally {
      setManagingPortal(false);
    }
  }, [demoSession]);

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <svg
          className="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
          />
        </svg>
        Back to demos
      </Link>

      {/* Active subscription card */}
      {demoSession && (
        <ActiveSubscriptionCard
          session={demoSession}
          onManage={handleManagePortal}
          managingPortal={managingPortal}
        />
      )}

      {/* Cancel banner */}
      {canceled && !demoSession && <CancelBanner />}

      {/* Header */}
      <div className="text-center max-w-2xl mx-auto mb-12">
        <span className="inline-block bg-primary/10 text-primary text-sm font-medium px-3 py-1 rounded-full mb-4">
          Live Stripe Demo
        </span>
        <h1 className="text-4xl sm:text-5xl font-bold text-foreground tracking-tight">
          Simple, transparent pricing
        </h1>
        <p className="text-lg text-muted-foreground mt-4">
          {demoSession
            ? "You can switch plans, upgrade, or manage your subscription below."
            : "Choose the plan that fits your needs. This demo uses real Stripe Checkout in test mode. No real charges will be made."}
        </p>
      </div>

      {/* Interval selector */}
      <div className="mb-10">
        <IntervalSelector value={interval} onChange={setInterval} />
      </div>

      {/* Error */}
      {error && (
        <div className="max-w-5xl mx-auto mb-6">
          <p className="text-sm text-destructive text-center">{error}</p>
        </div>
      )}

      {/* Pricing cards */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto mb-16">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="bg-card border border-border rounded-xl p-6 shadow-sm animate-pulse h-96"
            />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto mb-16">
          {tiers.map((tier) => (
            <PricingCard
              key={tier.tierName}
              tier={tier}
              interval={interval}
              onCheckout={handleCheckout}
              loading={checkoutLoading}
              currentPlanName={demoSession?.plan_name}
            />
          ))}
        </div>
      )}

      {/* Test card hint */}
      <div className="max-w-2xl mx-auto">
        <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
          <h3 className="text-sm font-semibold text-foreground mb-3">
            Stripe test mode — how to test
          </h3>
          <div className="space-y-2 text-sm text-muted-foreground">
            <p>
              This demo uses Stripe in{" "}
              <strong className="text-foreground">test mode</strong>. No real
              charges will be made. Use these test credentials:
            </p>
            <div className="bg-muted rounded-lg p-3 font-mono text-xs space-y-1">
              <p>
                <span className="text-foreground">Card number:</span> 4242 4242
                4242 4242
              </p>
              <p>
                <span className="text-foreground">Expiry:</span> Any future date
                (e.g. 12/28)
              </p>
              <p>
                <span className="text-foreground">CVC:</span> Any 3 digits (e.g.
                123)
              </p>
              <p>
                <span className="text-foreground">Name / ZIP:</span> Any values
              </p>
            </div>
            {demoSession && (
              <p className="text-xs">
                Use the <strong>Manage subscription</strong> button above to
                upgrade, downgrade, cancel, or update your payment method via the
                Stripe billing portal.
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
