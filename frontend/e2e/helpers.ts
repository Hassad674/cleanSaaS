import { type Page } from "@playwright/test";

export const API_URL = process.env.API_URL || "http://localhost:8081";

export async function login(page: Page, email: string, password: string) {
  await page.goto("/login");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /sign in|log in/i }).click();
  await page.waitForURL("**/dashboard");
}

export async function register(
  page: Page,
  name: string,
  email: string,
  password: string
) {
  await page.goto("/register");
  await page.getByLabel("Name").fill(name);
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /sign up|register|create/i }).click();
  await page.waitForURL("**/dashboard");
}
