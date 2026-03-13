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

export async function forgotPassword(email: string) {
  return api<{ message: string }>("/auth/forgot-password", {
    method: "POST",
    body: { email },
  });
}

export async function resetPassword(token: string, password: string) {
  return api<{ message: string }>("/auth/reset-password", {
    method: "POST",
    body: { token, password },
  });
}

export async function verifyEmail(token: string) {
  return api<{ message: string }>("/auth/verify-email", {
    method: "POST",
    body: { token },
  });
}

export async function resendVerification(authToken: string) {
  return api<{ message: string }>("/auth/resend-verification", {
    method: "POST",
    token: authToken,
  });
}
