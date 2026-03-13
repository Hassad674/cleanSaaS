"use server";

import { api } from "@/shared/lib/api";
import type {
  ReferralCodeResponse,
  ReferralStats,
  ReferralListResponse,
  ApplyReferralResponse,
} from "@/features/referral/types";

export async function getReferralCode(authToken: string) {
  return api<ReferralCodeResponse>("/referral/code", {
    token: authToken,
  });
}

export async function getReferralStats(authToken: string) {
  return api<ReferralStats>("/referral/stats", {
    token: authToken,
  });
}

export async function listReferrals(
  authToken: string,
  page: number = 1,
  limit: number = 20
) {
  const params = new URLSearchParams({
    page: String(page),
    limit: String(limit),
  });

  return api<ReferralListResponse>(`/referral/list?${params.toString()}`, {
    token: authToken,
  });
}

export async function applyReferral(authToken: string, code: string) {
  return api<ApplyReferralResponse>("/referral/apply", {
    method: "POST",
    token: authToken,
    body: { code },
  });
}
