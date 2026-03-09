const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

type FetchOptions = {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
  token?: string;
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
  const { method = "GET", body, headers = {}, token } = options;

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
    });

    const data = res.ok ? await res.json() : null;
    const error = res.ok ? null : (await res.json().catch(() => null))?.error || "Request failed";

    return { data, error, status: res.status };
  } catch {
    return { data: null, error: "Network error", status: 0 };
  }
}
