"use client";

import { useState, useEffect } from "react";
import { getPlans, createCheckout } from "@/features/billing/actions/billing";
import { formatCurrency } from "@/shared/lib/utils";
import { AUTH_TOKEN_KEY } from "@/shared/lib/constants";
import type { Plan } from "@/features/billing/types";

const POPULAR_PLAN_NAME = "Pro";

export function PricingCards() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [checkoutLoading, setCheckoutLoading] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getPlans().then((res) => {
      if (res.data) {
        setPlans(res.data);
      } else {
        setError(res.error ?? "Failed to load plans");
      }
      setLoading(false);
    });
  }, []);

  async function handleCheckout(planID: string) {
    const token =
      typeof window !== "undefined"
        ? localStorage.getItem(AUTH_TOKEN_KEY)
        : null;

    if (!token) {
      window.location.href = "/login";
      return;
    }

    setCheckoutLoading(planID);
    const res = await createCheckout(planID, token);
    setCheckoutLoading(null);

    if (res.data?.url) {
      window.location.href = res.data.url;
    } else {
      setError(res.error ?? "Failed to create checkout session");
    }
  }

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="bg-card border border-border rounded-xl p-6 shadow-sm animate-pulse h-80"
          />
        ))}
      </div>
    );
  }

  if (error && plans.length === 0) {
    return (
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm text-center">
        <p className="text-destructive">{error}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {error && (
        <p className="text-sm text-destructive text-center">{error}</p>
      )}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {plans.map((plan) => {
          const isPopular = plan.name === POPULAR_PLAN_NAME;
          const isFree = plan.price_cents === 0;

          return (
            <div
              key={plan.id}
              className={`bg-card border rounded-xl p-6 shadow-sm flex flex-col relative ${
                isPopular ? "border-primary ring-2 ring-primary" : "border-border"
              }`}
            >
              {isPopular && (
                <span className="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary text-primary-foreground text-xs font-medium px-3 py-1 rounded-full">
                  Most popular
                </span>
              )}

              <h3 className="text-lg font-semibold text-foreground">
                {plan.name}
              </h3>

              <div className="mt-4 mb-6">
                <span className="text-4xl font-bold text-foreground">
                  {isFree ? "Free" : formatCurrency(plan.price_cents, "usd")}
                </span>
                {!isFree && (
                  <span className="text-muted-foreground ml-1">
                    /{plan.interval === "month" ? "mo" : "yr"}
                  </span>
                )}
              </div>

              <ul className="space-y-3 mb-8 flex-1">
                {plan.features.map((feature) => (
                  <li
                    key={feature}
                    className="flex items-start gap-2 text-sm text-foreground"
                  >
                    <svg
                      className="h-5 w-5 shrink-0 text-primary mt-0.5"
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
                    {feature}
                  </li>
                ))}
              </ul>

              {isFree ? (
                <button
                  disabled
                  className="w-full bg-muted text-muted-foreground rounded-lg px-4 py-2.5 font-medium cursor-default"
                >
                  Current plan
                </button>
              ) : (
                <button
                  onClick={() => handleCheckout(plan.id)}
                  disabled={checkoutLoading === plan.id}
                  className="w-full bg-primary text-primary-foreground rounded-lg px-4 py-2.5 font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
                >
                  {checkoutLoading === plan.id ? "Redirecting..." : "Get started"}
                </button>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
