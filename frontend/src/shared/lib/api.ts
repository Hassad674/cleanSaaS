const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

type FetchOptions = {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
  token?: string;
  signal?: AbortSignal;
};

type ApiResponse<T> = {
  data: T | null;
  error: string | null;
  status: number;
};

export async function api<T>(
  path: string,
  options: FetchOptions = {}
): Promise<ApiResponse<T>> {
  const { method = "GET", body, headers = {}, token, signal } = options;

  const requestHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...headers,
  };

  if (token) {
    requestHeaders["Authorization"] = `Bearer ${token}`;
  }

  try {
    const res = await fetch(`${API_URL}${path}`, {
      method,
      headers: requestHeaders,
      body: body ? JSON.stringify(body) : undefined,
      signal,
    });

    // Parse body once, then decide if it's data or error
    const json = await res.json().catch(() => null);

    if (res.ok) {
      return { data: json as T, error: null, status: res.status };
    }

    const errorMessage = (json as { error?: string })?.error || "Request failed";
    return { data: null, error: errorMessage, status: res.status };
  } catch {
    return { data: null, error: "Network error", status: 0 };
  }
}
