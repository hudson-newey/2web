import { expect, test } from "playwright/test";

test("should have the correct title", async ({ page }) => {
  await page.goto("/");
  expect(await page.title()).toBe("Welcome to 2Web");
});
