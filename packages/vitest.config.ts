import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    dir: ".",
    include: ["**/*.spec.ts"],
    reporters: ["default"],
    environment: "happy-dom",
  },
  server: {
    port: 3000,
    open: false,
  },
});
