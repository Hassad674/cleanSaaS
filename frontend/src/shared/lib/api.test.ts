import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { api } from "./api";

// Save the original fetch
const originalFetch = global.fetch;

describe("api", () => {
  beforeEach(() => {
    // Reset fetch mock before each test
    global.fetch = vi.fn();
  });

  afterEach(() => {
    global.fetch = originalFetch;
  });

  it("makes GET request by default", async () => {
    const mockData = { id: 1, name: "test" };
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => mockData,
    });

    const result = await api("/test");

    expect(result.data).toEqual(mockData);
    expect(result.error).toBeNull();
    expect(result.status).toBe(200);

    const [url, options] = (global.fetch as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(url).toContain("/test");
    expect(options.method).toBe("GET");
    expect(options.headers["Content-Type"]).toBe("application/json");
  });

  it("sends POST request with body", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 201,
      json: async () => ({ created: true }),
    });

    const body = { name: "new item" };
    const result = await api("/items", { method: "POST", body });

    expect(result.data).toEqual({ created: true });
    expect(result.status).toBe(201);

    const [, options] = (global.fetch as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(options.method).toBe("POST");
    expect(options.body).toBe(JSON.stringify(body));
  });

  it("includes authorization header when token provided", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({}),
    });

    await api("/protected", { token: "my-jwt-token" });

    const [, options] = (global.fetch as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(options.headers["Authorization"]).toBe("Bearer my-jwt-token");
  });

  it("does not include authorization header when no token", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({}),
    });

    await api("/public");

    const [, options] = (global.fetch as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(options.headers["Authorization"]).toBeUndefined();
  });

  it("returns error for non-ok response", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 404,
      json: async () => ({ error: "not found" }),
    });

    const result = await api("/missing");

    expect(result.data).toBeNull();
    expect(result.error).toBe("not found");
    expect(result.status).toBe(404);
  });

  it("returns fallback error when error response has no JSON", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error("no json");
      },
    });

    const result = await api("/broken");

    expect(result.data).toBeNull();
    expect(result.error).toBe("Request failed");
    expect(result.status).toBe(500);
  });

  it("handles network error", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValueOnce(
      new Error("Network failure")
    );

    const result = await api("/offline");

    expect(result.data).toBeNull();
    expect(result.error).toBe("Network error");
    expect(result.status).toBe(0);
  });

  it("does not send body for GET requests", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({}),
    });

    await api("/items");

    const [, options] = (global.fetch as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(options.body).toBeUndefined();
  });

  it("merges custom headers", async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({}),
    });

    await api("/test", { headers: { "X-Custom": "value" } });

    const [, options] = (global.fetch as ReturnType<typeof vi.fn>).mock
      .calls[0];
    expect(options.headers["X-Custom"]).toBe("value");
    expect(options.headers["Content-Type"]).toBe("application/json");
  });
});
