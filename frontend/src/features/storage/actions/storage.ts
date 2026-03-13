"use server";

import { api } from "@/shared/lib/api";
import type { FileItem, FilesResponse } from "@/features/storage/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

export async function uploadFile(formData: FormData, authToken: string) {
  try {
    const res = await fetch(`${API_URL}/files/upload`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      body: formData,
    });

    const data = res.ok ? await res.json() : null;
    const error = res.ok
      ? null
      : (await res.json().catch(() => null))?.error || "Upload failed";

    return { data: data as FileItem | null, error, status: res.status };
  } catch {
    return { data: null, error: "Network error", status: 0 };
  }
}

export async function getFiles(
  authToken: string,
  page: number = 1,
  limit: number = 20
) {
  return api<FilesResponse>(`/files?page=${page}&limit=${limit}`, {
    token: authToken,
  });
}

export async function deleteFile(fileId: string, authToken: string) {
  return api<{ message: string }>(`/files/${fileId}`, {
    method: "DELETE",
    token: authToken,
  });
}
