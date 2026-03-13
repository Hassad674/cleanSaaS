"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/features/auth/hooks/use-auth";
import {
  getSubscription,
  getPlans,
  cancelSubscription,
  createPortalSession,
} from "@/features/billing/actions/billing";
import { formatDate } from "@/shared/lib/utils";
import type { Subscription, Plan } from "@/features/billing/types";

const STATUS_STYLES: Record<string, string> = {
  active: "bg-accent text-primary",
  trialing: "bg-accent text-primary",
  canceled: "bg-muted text-muted-foreground",
  past_due: "bg-destructive/10 text-destructive",
  inactive: "bg-muted text-muted-foreground",
};

export function SubscriptionStatus() {
  const { getToken } = useAuth({ required: true });

  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    Promise.all([getSubscription(token), getPlans()]).then(
      ([subRes, plansRes]) => {
        if (subRes.data) setSubscription(subRes.data);
        if (plansRes.data) setPlans(plansRes.data);
        setLoading(false);
      }
    );
  }, [getToken]);

  const currentPlan = plans.find((p) => p.id === subscription?.plan_id);

  async function handleCancel() {
    const token = getToken();
    if (!token) return;

    setActionLoading(true);
    setError(null);
    const res = await cancelSubscription(token);
    setActionLoading(false);
    setShowConfirm(false);

    if (res.error) {
      setError(res.error);
    } else {
      setSuccess("Subscription will be canceled at the end of the billing period.");
      setSubscription((prev) =>
        prev ? { ...prev, cancel_at_period_end: true } : prev
      );
    }
  }

  async function handlePortal() {
    const token = getToken();
    if (!token) return;

    setActionLoading(true);
    setError(null);
    const res = await createPortalSession(token);
    setActionLoading(false);

    if (res.data?.url) {
      window.location.href = res.data.url;
    } else {
      setError(res.error ?? "Failed to open billing portal");
    }
  }

  if (loading) {
    return (
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm animate-pulse h-40" />
    );
  }

  if (!subscription) {
    return (
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-foreground mb-2">
          Subscription
        </h2>
        <p className="text-muted-foreground">No active subscription.</p>
      </div>
    );
  }

  const statusLabel = subscription.status.replace("_", " ");
  const statusClass = STATUS_STYLES[subscription.status] ?? STATUS_STYLES.inactive;

  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">Subscription</h2>
        <span
          className={`text-xs font-medium px-2.5 py-0.5 rounded-full capitalize ${statusClass}`}
        >
          {statusLabel}
        </span>
      </div>

      <div className="space-y-1 text-sm">
        <p className="text-foreground">
          <span className="text-muted-foreground">Plan:</span>{" "}
          {currentPlan?.name ?? "Unknown"}
        </p>
        <p className="text-foreground">
          <span className="text-muted-foreground">Next billing date:</span>{" "}
          {formatDate(subscription.current_period_end)}
        </p>
        {subscription.cancel_at_period_end && (
          <p className="text-sm text-muted-foreground">
            Your subscription will end on{" "}
            {formatDate(subscription.current_period_end)}.
          </p>
        )}
      </div>

      {success && <p className="text-sm text-primary">{success}</p>}
      {error && <p className="text-sm text-destructive">{error}</p>}

      <div className="flex gap-3 pt-2">
        <button
          onClick={handlePortal}
          disabled={actionLoading}
          className="bg-primary text-primary-foreground rounded-lg px-4 py-2 text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          Manage billing
        </button>

        {subscription.status === "active" && !subscription.cancel_at_period_end && (
          <>
            {showConfirm ? (
              <div className="flex items-center gap-2">
                <span className="text-sm text-muted-foreground">
                  Are you sure?
                </span>
                <button
                  onClick={handleCancel}
                  disabled={actionLoading}
                  className="bg-destructive text-primary-foreground rounded-lg px-3 py-2 text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
                >
                  {actionLoading ? "Canceling..." : "Yes, cancel"}
                </button>
                <button
                  onClick={() => setShowConfirm(false)}
                  className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  No, keep it
                </button>
              </div>
            ) : (
              <button
                onClick={() => setShowConfirm(true)}
                className="border border-border text-foreground rounded-lg px-4 py-2 text-sm font-medium hover:bg-muted transition-colors"
              >
                Cancel subscription
              </button>
            )}
          </>
        )}
      </div>
    </div>
  );
}
