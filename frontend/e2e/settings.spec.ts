import { test, expect } from "@playwright/test";
import { login } from "./helpers";

test.describe("Settings", () => {
  test.beforeEach(async ({ page }) => {
    await login(page, "admin@cleansaas.com", "admin123");
  });

  test("navigate to settings page", async ({ page }) => {
    await page.goto("/settings");
    await expect(page.getByText(/settings|profile/i)).toBeVisible();
  });

  test("update user name", async ({ page }) => {
    await page.goto("/settings");
    const nameInput = page.getByLabel("Name");
    await nameInput.clear();
    await nameInput.fill("Updated Admin");
    await page.getByRole("button", { name: /save|update/i }).click();
    await expect(page.getByText(/updated|saved|success/i)).toBeVisible();
  });
});
