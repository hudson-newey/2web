import { Page } from "@playwright/test";

const serverLocation = "http://localhost:3000";

export async function navigateToPage(page: Page, path: string) {
  await page.goto(serverLocation + path)
  await page.waitForLoadState("networkidle");
}
