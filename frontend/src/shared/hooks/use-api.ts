"use client";

import { useState, useCallback } from "react";
import { api } from "@/shared/lib/api";

type UseApiState<T> = {
  data: T | null;
  error: string | null;
  loading: boolean;
};

export function useApi<T>() {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    error: null,
    loading: false,
  });

  const execute = useCallback(
    async (path: string, options?: Parameters<typeof api>[1]) => {
      setState((prev) => ({ ...prev, loading: true, error: null }));

      const result = await api<T>(path, options);

      setState({
        data: result.data,
        error: result.error,
        loading: false,
      });

      return result;
    },
    []
  );

  return { ...state, execute };
}
