import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 1,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    trace: "on-first-retry",
    bypassCSP: true,
    screenshot: "only-on-failure",
  },
  projects: [{ name: "chromium" }, { name: "firefox" }, { name: "webkit" }],
  webServer: {
    // we run the dev server on port 3000 so that we can run a separate
    // development server that is independent of the test server
    command: "pnpm dev --port 3000",
    port: 3000,
    reuseExistingServer: true,
  },
});
