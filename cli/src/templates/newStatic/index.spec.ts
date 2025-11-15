import { expect, test } from "playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("/");
});

test("should have the correct title", async ({ page }) => {
  expect(await page.title()).toBe("2Web App");
});

test.describe("counter", () => {
  test("should increment and decrement the counter", async ({ page }) => {
    const counter = page.locator(".counter-btn");

    // Assert that the counter starts at 0
    await expect(counter).toHaveText("0");

    // Assert that we can increment the counter by clicking on the button
    await counter.click();
    await expect(counter).toHaveText("1");
  });
});
