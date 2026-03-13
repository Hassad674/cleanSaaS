"use server";

import { api } from "@/shared/lib/api";
import type { Plan, Subscription, Invoice } from "@/features/billing/types";

export async function getPlans() {
  return api<Plan[]>("/billing/plans");
}

export async function createCheckout(planID: string, token: string) {
  return api<{ url: string }>("/billing/checkout", {
    method: "POST",
    body: { plan_id: planID },
    token,
  });
}

export async function getSubscription(token: string) {
  return api<Subscription>("/billing/subscription", { token });
}

export async function cancelSubscription(token: string) {
  return api<{ message: string }>("/billing/cancel", {
    method: "POST",
    token,
  });
}

export async function createPortalSession(token: string) {
  return api<{ url: string }>("/billing/portal", {
    method: "POST",
    token,
  });
}

export async function getInvoices(token: string, offset: number, limit: number) {
  return api<{ invoices: Invoice[]; total: number }>(
    `/billing/invoices?offset=${offset}&limit=${limit}`,
    { token }
  );
}
