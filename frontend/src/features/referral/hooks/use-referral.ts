"use client";

import { useState, useEffect, useCallback } from "react";
import { useAuth } from "@/shared/hooks/use-auth";
import {
  getReferralCode as getReferralCodeAction,
  getReferralStats as getReferralStatsAction,
  listReferrals as listReferralsAction,
} from "@/features/referral/actions/referral";
import { PAGINATION_DEFAULT_LIMIT } from "@/shared/lib/constants";
import type { Referral, ReferralStats } from "@/features/referral/types";

export function useReferral() {
  const { getToken } = useAuth({ required: true });

  const [code, setCode] = useState<string | null>(null);
  const [stats, setStats] = useState<ReferralStats | null>(null);
  const [referrals, setReferrals] = useState<Referral[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const limit = PAGINATION_DEFAULT_LIMIT;

  const fetchCode = useCallback(() => {
    const token = getToken();
    if (!token) return;

    getReferralCodeAction(token).then((res) => {
      if (res.data) {
        setCode(res.data.code);
      } else {
        setError(res.error ?? "Failed to load referral code");
      }
    });
  }, [getToken]);

  const fetchStats = useCallback(() => {
    const token = getToken();
    if (!token) return;

    getReferralStatsAction(token).then((res) => {
      if (res.data) {
        setStats(res.data);
      }
    });
  }, [getToken]);

  const fetchReferrals = useCallback(() => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);
    listReferralsAction(token, page, limit).then((res) => {
      if (res.data) {
        setReferrals(res.data.referrals ?? []);
        setTotal(res.data.total);
      } else {
        setError(res.error ?? "Failed to load referrals");
      }
      setLoading(false);
    });
  }, [getToken, page, limit]);

  // Fetch code and stats on mount
  useEffect(() => {
    fetchCode();
    fetchStats();
  }, [fetchCode, fetchStats]);

  // Fetch referrals when page changes
  useEffect(() => {
    fetchReferrals();
  }, [fetchReferrals]);

  const totalPages = Math.ceil(total / limit);
  const hasNext = page < totalPages;
  const hasPrev = page > 1;

  const goToNextPage = useCallback(() => {
    if (hasNext) setPage((prev) => prev + 1);
  }, [hasNext]);

  const goToPrevPage = useCallback(() => {
    if (hasPrev) setPage((prev) => prev - 1);
  }, [hasPrev]);

  return {
    code,
    stats,
    referrals,
    total,
    page,
    totalPages,
    loading,
    error,
    hasNext,
    hasPrev,
    goToNextPage,
    goToPrevPage,
    fetchCode,
    fetchStats,
    fetchReferrals,
  };
}
