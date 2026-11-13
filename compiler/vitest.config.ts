import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default defineConfig(() =>
  mergeConfig(
    viteConfig,
    defineConfig({
      test: {
        dir: ".",
        include: ["**/*.spec.ts"],
        forceRerunTriggers: ["/dist/**"],
        setupFiles: ["./tests/vitest.setup.ts"],
        reporters: ["default"],
        environment: "happy-dom",
        globals: true,
      },
    }),
  ),
);
