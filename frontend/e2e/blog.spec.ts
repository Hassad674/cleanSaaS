import { test, expect } from "@playwright/test";

test.describe("Blog", () => {
  test("blog listing page loads", async ({ page }) => {
    await page.goto("/blog");
    await expect(page.getByText(/blog/i)).toBeVisible();
  });

  test("blog post page displays content", async ({ page }) => {
    await page.goto("/blog");

    // Click the first blog post link
    const firstPost = page.getByRole("link").filter({ hasText: /.+/ }).first();
    const postTitle = await firstPost.textContent();

    if (postTitle) {
      await firstPost.click();
      await expect(page.getByText(postTitle)).toBeVisible();
    }
  });
});
