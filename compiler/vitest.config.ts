import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    dir: ".",
    include: ["**/*.spec.ts"],
    forceRerunTriggers: ["**/dist/**"],
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
