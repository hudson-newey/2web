import { defineConfig } from "vite";
import dts from "vite-plugin-dts";

export default defineConfig({
  appType: "mpa",
  plugins: [dts()],
  build: {
    outDir: "./dist",
    copyPublicDir: false,
    rollupOptions: {
      external: ["vite", "happy-dom", "express", "node", /node:/],
    },
    lib: {
      // prettier-ignore
      entry: {
        "index": "./index.ts",
        "animations": "./animations/index.ts",
        "database": "./database/index.ts",
        "pre-fetcher": "./pre-fetcher/index.ts",
        "route-guards": "./route-guards/index.ts",
        "signals": "./signals/index.ts",
        "ssr": "./ssr/index.ts",
        "threads": "./threads/index.ts",
      },
      // Since many of the packages are framework agnostic, I want to distribute
      // the compiled JavaScript in both esm and cjs so that it can be consumed
      // by as many targets as possible.
      formats: ["es", "cjs"],
    },
  },
});
