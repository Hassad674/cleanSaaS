"use server";

import { api } from "@/shared/lib/api";
import type { AuthResponse } from "@/shared/types/common";

export async function login(email: string, password: string) {
  return api<AuthResponse>("/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export async function register(email: string, name: string, password: string) {
  return api<AuthResponse>("/auth/register", {
    method: "POST",
    body: { email, name, password },
  });
}
