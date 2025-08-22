import { defineConfig } from "vite";
import dts from "vite-plugin-dts";

export default defineConfig({
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
        // Note: there is no barrel file to import all two-web/kit packages
        // I do this so that we get improved tree shaking.
        "database": "./database/index.ts",
        "pre-fetcher": "./pre-fetcher/index.ts",
        "route-guards": "./route-guards/index.ts",
        "signals": "./signals/index.ts",
        "ssr": "./ssr/index.ts",
        // "vite-plugin": "./vite-plugin/index.ts",
      },
      // Since many of the packages are framework agnostic, I want to distribute
      // the compiled JavaScript in as many formats as possible so that it can
      // be consumed by any framework.
      formats: ["es"],
    },
  },
});
