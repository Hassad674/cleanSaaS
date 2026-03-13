import { test, expect } from "@playwright/test";
import { register, login } from "./helpers";

const TEST_USER = {
  name: "Test User",
  email: `test-${Date.now()}@example.com`,
  password: "SecurePass123!",
};

test.describe("Authentication", () => {
  test("register new user and see dashboard", async ({ page }) => {
    await register(page, TEST_USER.name, TEST_USER.email, TEST_USER.password);
    await expect(page).toHaveURL(/.*dashboard/);
    await expect(page.getByText(/dashboard/i)).toBeVisible();
  });

  test("login with existing user", async ({ page }) => {
    await login(page, "admin@cleansaas.com", "admin123");
    await expect(page).toHaveURL(/.*dashboard/);
  });

  test("login page shows error for wrong password", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel("Email").fill("admin@cleansaas.com");
    await page.getByLabel("Password").fill("wrongpassword");
    await page.getByRole("button", { name: /sign in|log in/i }).click();
    await expect(page.getByText(/invalid|incorrect|error/i)).toBeVisible();
  });
});
