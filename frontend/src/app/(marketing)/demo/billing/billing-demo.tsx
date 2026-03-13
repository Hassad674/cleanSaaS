"use client";

import { useState, useCallback } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

/* -------------------------------------------------------------------------- */
/*  Mock data                                                                 */
/* -------------------------------------------------------------------------- */

type PlanFeature = string;

interface Plan {
  id: string;
  name: string;
  monthlyPrice: number;
  annualPrice: number;
  features: PlanFeature[];
  highlighted: boolean;
  cta: string;
}

const plans: Plan[] = [
  {
    id: "free",
    name: "Free",
    monthlyPrice: 0,
    annualPrice: 0,
    features: [
      "1 project",
      "100MB storage",
      "Community support",
      "Basic analytics",
    ],
    highlighted: false,
    cta: "Get started",
  },
  {
    id: "pro",
    name: "Pro",
    monthlyPrice: 19,
    annualPrice: 15,
    features: [
      "Unlimited projects",
      "10GB storage",
      "Priority support",
      "Advanced analytics",
      "AI assistant",
      "Custom domains",
    ],
    highlighted: true,
    cta: "Get started",
  },
  {
    id: "enterprise",
    name: "Enterprise",
    monthlyPrice: 49,
    annualPrice: 39,
    features: [
      "Everything in Pro",
      "Unlimited storage",
      "Dedicated support",
      "SSO & SAML",
      "Custom integrations",
      "SLA guarantee",
    ],
    highlighted: false,
    cta: "Get started",
  },
];

interface MockInvoice {
  id: string;
  date: string;
  plan: string;
  amount: string;
  status: "paid" | "pending";
}

const mockInvoices: MockInvoice[] = [
  {
    id: "inv-001",
    date: "Mar 1, 2026",
    plan: "Pro",
    amount: "$19.00",
    status: "paid",
  },
  {
    id: "inv-002",
    date: "Feb 1, 2026",
    plan: "Pro",
    amount: "$19.00",
    status: "paid",
  },
  {
    id: "inv-003",
    date: "Jan 1, 2026",
    plan: "Pro",
    amount: "$19.00",
    status: "paid",
  },
  {
    id: "inv-004",
    date: "Dec 1, 2025",
    plan: "Free",
    amount: "$0.00",
    status: "pending",
  },
];

/* -------------------------------------------------------------------------- */
/*  Check icon                                                                */
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
/*  Billing toggle                                                            */
/* -------------------------------------------------------------------------- */

function BillingToggle({
  annual,
  onChange,
}: {
  annual: boolean;
  onChange: (annual: boolean) => void;
}) {
  return (
    <div className="flex items-center justify-center gap-3">
      <span
        className={cn(
          "text-sm font-medium transition-colors",
          !annual ? "text-foreground" : "text-muted-foreground"
        )}
      >
        Monthly
      </span>
      <button
        type="button"
        role="switch"
        aria-checked={annual}
        onClick={() => onChange(!annual)}
        className={cn(
          "relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
          annual ? "bg-primary" : "bg-muted"
        )}
      >
        <span
          className={cn(
            "pointer-events-none inline-block h-5 w-5 rounded-full bg-card shadow-sm ring-0 transition-transform",
            annual ? "translate-x-5" : "translate-x-0"
          )}
        />
      </button>
      <span
        className={cn(
          "text-sm font-medium transition-colors",
          annual ? "text-foreground" : "text-muted-foreground"
        )}
      >
        Annual
      </span>
      {annual && (
        <span className="text-xs font-medium text-primary bg-primary/10 px-2 py-0.5 rounded-full">
          Save 20%
        </span>
      )}
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Pricing card                                                              */
/* -------------------------------------------------------------------------- */

function PricingCard({
  plan,
  annual,
  onSelect,
}: {
  plan: Plan;
  annual: boolean;
  onSelect: (plan: Plan) => void;
}) {
  const price = annual ? plan.annualPrice : plan.monthlyPrice;
  const isFree = price === 0;

  return (
    <div
      className={cn(
        "bg-card border rounded-xl p-6 shadow-sm flex flex-col relative",
        plan.highlighted
          ? "border-primary ring-2 ring-primary"
          : "border-border"
      )}
    >
      {plan.highlighted && (
        <span className="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary text-primary-foreground text-xs font-medium px-3 py-1 rounded-full">
          Most popular
        </span>
      )}

      <h3 className="text-lg font-semibold text-foreground">{plan.name}</h3>

      <div className="mt-4 mb-6">
        <span className="text-4xl font-bold text-foreground">
          {isFree ? "Free" : `$${price}`}
        </span>
        {!isFree && (
          <span className="text-muted-foreground ml-1">/mo</span>
        )}
        {!isFree && annual && (
          <span className="block text-xs text-muted-foreground mt-1">
            Billed annually (${price * 12}/yr)
          </span>
        )}
      </div>

      <ul className="space-y-3 mb-8 flex-1">
        {plan.features.map((feature) => (
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
          onClick={() => onSelect(plan)}
          className="w-full bg-muted text-foreground rounded-lg px-4 py-2.5 font-medium hover:bg-accent transition-colors"
        >
          {plan.cta}
        </button>
      ) : (
        <button
          onClick={() => onSelect(plan)}
          className="w-full bg-primary text-primary-foreground rounded-lg px-4 py-2.5 font-medium hover:opacity-90 transition-opacity"
        >
          {plan.cta}
        </button>
      )}
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Checkout modal                                                            */
/* -------------------------------------------------------------------------- */

type CheckoutStep = "form" | "processing" | "success";

function CheckoutModal({
  plan,
  annual,
  onClose,
}: {
  plan: Plan;
  annual: boolean;
  onClose: () => void;
}) {
  const [step, setStep] = useState<CheckoutStep>("form");
  const price = annual ? plan.annualPrice : plan.monthlyPrice;

  const handlePay = useCallback(() => {
    setStep("processing");
    const timer = setTimeout(() => {
      setStep("success");
    }, 2000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-foreground/40"
        onClick={step !== "processing" ? onClose : undefined}
        aria-hidden="true"
      />

      {/* Modal */}
      <div className="relative bg-card border border-border rounded-xl shadow-lg w-full max-w-md p-6 space-y-6">
        {step === "form" && (
          <>
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">
                Checkout
              </h2>
              <button
                onClick={onClose}
                className="text-muted-foreground hover:text-foreground transition-colors"
                aria-label="Close"
              >
                <svg
                  className="h-5 w-5"
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
              </button>
            </div>

            {/* Order summary */}
            <div className="bg-muted rounded-lg p-4">
              <div className="flex items-center justify-between text-sm">
                <span className="text-foreground font-medium">
                  {plan.name} plan
                </span>
                <span className="text-foreground font-semibold">
                  ${price}/mo
                </span>
              </div>
              {annual && (
                <p className="text-xs text-muted-foreground mt-1">
                  Billed annually at ${price * 12}/yr
                </p>
              )}
            </div>

            {/* Card form */}
            <div className="space-y-4">
              <div>
                <label
                  htmlFor="card-number"
                  className="block text-sm font-medium text-foreground mb-1.5"
                >
                  Card number
                </label>
                <input
                  id="card-number"
                  type="text"
                  defaultValue="4242 4242 4242 4242"
                  readOnly
                  className="w-full bg-muted border border-border rounded-lg px-3 py-2 text-sm text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label
                    htmlFor="card-expiry"
                    className="block text-sm font-medium text-foreground mb-1.5"
                  >
                    Expiry
                  </label>
                  <input
                    id="card-expiry"
                    type="text"
                    defaultValue="12/28"
                    readOnly
                    className="w-full bg-muted border border-border rounded-lg px-3 py-2 text-sm text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  />
                </div>
                <div>
                  <label
                    htmlFor="card-cvc"
                    className="block text-sm font-medium text-foreground mb-1.5"
                  >
                    CVC
                  </label>
                  <input
                    id="card-cvc"
                    type="text"
                    defaultValue="123"
                    readOnly
                    className="w-full bg-muted border border-border rounded-lg px-3 py-2 text-sm text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  />
                </div>
              </div>
            </div>

            <button
              onClick={handlePay}
              className="w-full bg-primary text-primary-foreground rounded-lg px-4 py-2.5 font-medium hover:opacity-90 transition-opacity"
            >
              Pay ${price.toFixed(2)}
            </button>

            <p className="text-xs text-muted-foreground text-center">
              This is a simulated checkout. No real payment will be processed.
            </p>
          </>
        )}

        {step === "processing" && (
          <div className="flex flex-col items-center justify-center py-8 space-y-4">
            {/* Spinner */}
            <div className="h-10 w-10 rounded-full border-4 border-muted border-t-primary animate-spin" />
            <p className="text-sm text-muted-foreground font-medium">
              Processing payment...
            </p>
          </div>
        )}

        {step === "success" && (
          <div className="flex flex-col items-center justify-center py-8 space-y-4">
            <div className="h-14 w-14 rounded-full bg-primary/10 flex items-center justify-center">
              <CheckIcon className="h-8 w-8 text-primary" />
            </div>
            <div className="text-center space-y-1">
              <h3 className="text-lg font-semibold text-foreground">
                Payment successful!
              </h3>
              <p className="text-sm text-muted-foreground">
                You are now subscribed to the {plan.name} plan.
              </p>
            </div>
            <button
              onClick={onClose}
              className="bg-primary text-primary-foreground rounded-lg px-6 py-2 text-sm font-medium hover:opacity-90 transition-opacity"
            >
              Done
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Invoice table                                                             */
/* -------------------------------------------------------------------------- */

function InvoiceTable() {
  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm space-y-4">
      <h2 className="text-lg font-semibold text-foreground">Invoice history</h2>

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left">
              <th className="pb-3 font-medium text-muted-foreground">Date</th>
              <th className="pb-3 font-medium text-muted-foreground">Plan</th>
              <th className="pb-3 font-medium text-muted-foreground">Amount</th>
              <th className="pb-3 font-medium text-muted-foreground">Status</th>
            </tr>
          </thead>
          <tbody>
            {mockInvoices.map((invoice, index) => (
              <tr
                key={invoice.id}
                className={cn(
                  "border-b border-border last:border-0",
                  index % 2 === 1 && "bg-muted/50"
                )}
              >
                <td className="py-3 text-foreground">{invoice.date}</td>
                <td className="py-3 text-foreground">{invoice.plan}</td>
                <td className="py-3 text-foreground font-medium">
                  {invoice.amount}
                </td>
                <td className="py-3">
                  <span
                    className={cn(
                      "text-xs font-medium px-2.5 py-0.5 rounded-full capitalize",
                      invoice.status === "paid"
                        ? "bg-primary/10 text-primary"
                        : "bg-muted text-muted-foreground"
                    )}
                  >
                    {invoice.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Main demo component                                                       */
/* -------------------------------------------------------------------------- */

export function BillingDemo() {
  const [annual, setAnnual] = useState(false);
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null);

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

      {/* Header */}
      <div className="text-center max-w-2xl mx-auto mb-12">
        <span className="inline-block bg-primary/10 text-primary text-sm font-medium px-3 py-1 rounded-full mb-4">
          Billing Demo
        </span>
        <h1 className="text-4xl sm:text-5xl font-bold text-foreground tracking-tight">
          Simple, transparent pricing
        </h1>
        <p className="text-lg text-muted-foreground mt-4">
          Choose the plan that fits your needs. Upgrade, downgrade, or cancel at
          any time.
        </p>
      </div>

      {/* Billing toggle */}
      <div className="mb-10">
        <BillingToggle annual={annual} onChange={setAnnual} />
      </div>

      {/* Pricing cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto mb-16">
        {plans.map((plan) => (
          <PricingCard
            key={plan.id}
            plan={plan}
            annual={annual}
            onSelect={setSelectedPlan}
          />
        ))}
      </div>

      {/* Invoice table */}
      <div className="max-w-3xl mx-auto">
        <InvoiceTable />
      </div>

      {/* Disclaimer */}
      <p className="text-center text-sm text-muted-foreground mt-8">
        This is a simulated demo. In production, billing connects to Stripe via
        the Go backend.
      </p>

      {/* Checkout modal */}
      {selectedPlan && (
        <CheckoutModal
          plan={selectedPlan}
          annual={annual}
          onClose={() => setSelectedPlan(null)}
        />
      )}
    </div>
  );
}
