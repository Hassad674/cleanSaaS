"use server";

import { api } from "@/shared/lib/api";
import type { User } from "@/shared/types/common";
import type { UpdateProfileData, ChangePasswordData } from "@/features/user/types";

export async function getProfile(token: string) {
  return api<User>("/users/me", { token });
}

export async function updateProfile(token: string, data: UpdateProfileData) {
  return api<User>("/users/me", {
    method: "PATCH",
    body: data,
    token,
  });
}

export async function changePassword(token: string, data: ChangePasswordData) {
  return api<{ message: string }>("/users/me/password", {
    method: "PUT",
    body: data,
    token,
  });
}

export async function deleteAccount(token: string) {
  return api<void>("/users/me", {
    method: "DELETE",
    token,
  });
}
