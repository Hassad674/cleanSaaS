import { describe, it, expect } from "vitest";
import { cn, formatDate, formatCurrency } from "./utils";

describe("cn", () => {
  it("merges class names", () => {
    expect(cn("px-4", "py-2")).toBe("px-4 py-2");
  });

  it("handles conditional classes", () => {
    const isActive = true;
    expect(cn("base", isActive && "active")).toBe("base active");
  });

  it("handles false conditionals", () => {
    const isActive = false;
    expect(cn("base", isActive && "active")).toBe("base");
  });

  it("merges conflicting Tailwind classes (last wins)", () => {
    expect(cn("px-4", "px-6")).toBe("px-6");
  });

  it("handles empty input", () => {
    expect(cn()).toBe("");
  });

  it("handles undefined and null", () => {
    expect(cn("base", undefined, null, "extra")).toBe("base extra");
  });

  it("handles arrays of classes", () => {
    expect(cn(["px-4", "py-2"])).toBe("px-4 py-2");
  });
});

describe("formatDate", () => {
  it("formats a date string", () => {
    // Use a fixed date to avoid timezone issues
    const result = formatDate("2024-01-15T00:00:00Z");
    expect(result).toContain("Jan");
    expect(result).toContain("2024");
    expect(result).toContain("15");
  });

  it("formats a Date object", () => {
    const date = new Date("2024-06-01T00:00:00Z");
    const result = formatDate(date);
    expect(result).toContain("2024");
  });

  it("formats a recent date", () => {
    const result = formatDate("2025-12-25T12:00:00Z");
    expect(result).toContain("Dec");
    expect(result).toContain("25");
    expect(result).toContain("2025");
  });
});

describe("formatCurrency", () => {
  it("formats cents to dollars (USD)", () => {
    expect(formatCurrency(1900)).toBe("$19.00");
  });

  it("formats zero cents", () => {
    expect(formatCurrency(0)).toBe("$0.00");
  });

  it("formats small amounts", () => {
    expect(formatCurrency(99)).toBe("$0.99");
  });

  it("formats large amounts", () => {
    const result = formatCurrency(999900);
    expect(result).toContain("9,999.00");
  });

  it("formats with different currency", () => {
    const result = formatCurrency(1500, "eur");
    // EUR formatting varies by locale, just check it contains amount
    expect(result).toContain("15.00");
  });

  it("formats with explicit usd currency", () => {
    expect(formatCurrency(500, "usd")).toBe("$5.00");
  });
});
