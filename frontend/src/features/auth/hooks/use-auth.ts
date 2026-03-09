"use client";

import { useState, useEffect, useCallback } from "react";
import { AUTH_TOKEN_KEY } from "@/shared/lib/constants";
import type { User } from "@/shared/types/common";
import { api } from "@/shared/lib/api";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const getToken = useCallback(() => {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(AUTH_TOKEN_KEY);
  }, []);

  const setToken = useCallback((token: string) => {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    setUser(null);
  }, []);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      setLoading(false);
      return;
    }

    api<User>("/users/me", { token }).then((res) => {
      if (res.data) setUser(res.data);
      setLoading(false);
    });
  }, [getToken]);

  return { user, loading, getToken, setToken, logout };
}
