import { defineConfig } from "playwright/test";

// For more information about Playwright configuration, see:
// https://playwright.dev/docs/test-configuration
export default defineConfig({
  webServer: {
    command: "2web serve",
    port: 2000,
  },
  fullyParallel: true,
  reporter: [["list"]],
});
