import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    dir: "./tests",
    include: ["**/*.spec.ts", "**/*.test.ts"],
    globals: true,
    pool: "threads",
    maxConcurrency: process.env.CI ? 1 : undefined,
    retry: 1,
    reporters: ["html", "default"],
    outputFile: {
      html: "./html-report/index.html",
    },
    environment: "jsdom",
  },
  server: {
    port: 3000,
    open: false,
    watch: {
      usePolling: true,
    },
  },
});
