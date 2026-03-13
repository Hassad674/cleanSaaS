"use client";

import { useState, useEffect, useCallback } from "react";
import { useAuth } from "@/shared/hooks/use-auth";
import {
  getFiles,
  uploadFile as uploadFileAction,
  deleteFile as deleteFileAction,
} from "@/features/storage/actions/storage";
import { PAGINATION_DEFAULT_LIMIT } from "@/shared/lib/constants";
import type { FileItem } from "@/features/storage/types";

export function useStorage() {
  const { getToken } = useAuth({ required: true });

  const [files, setFiles] = useState<FileItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const limit = PAGINATION_DEFAULT_LIMIT;

  const fetchFiles = useCallback(() => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    setError(null);
    getFiles(token, page, limit).then((res) => {
      if (res.data) {
        setFiles(res.data.files ?? []);
        setTotal(res.data.total);
      } else {
        setError(res.error ?? "Failed to load files");
      }
      setLoading(false);
    });
  }, [getToken, page, limit]);

  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  const uploadFile = useCallback(
    async (file: File) => {
      const token = getToken();
      if (!token) return { error: "Not authenticated" };

      setUploading(true);
      setError(null);

      const formData = new FormData();
      formData.append("file", file);

      const res = await uploadFileAction(formData, token);
      setUploading(false);

      if (res.error) {
        setError(res.error);
        return { error: res.error };
      }

      // Refresh the file list after successful upload
      fetchFiles();
      return { error: null };
    },
    [getToken, fetchFiles]
  );

  const deleteFile = useCallback(
    async (fileId: string) => {
      const token = getToken();
      if (!token) return { error: "Not authenticated" };

      setError(null);
      const res = await deleteFileAction(fileId, token);

      if (res.error) {
        setError(res.error);
        return { error: res.error };
      }

      // Refresh the file list after successful deletion
      fetchFiles();
      return { error: null };
    },
    [getToken, fetchFiles]
  );

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
    files,
    total,
    page,
    totalPages,
    loading,
    uploading,
    error,
    hasNext,
    hasPrev,
    uploadFile,
    deleteFile,
    goToNextPage,
    goToPrevPage,
    refetch: fetchFiles,
  };
}
