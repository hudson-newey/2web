import { defineConfig } from "vite";
import dts from "vite-plugin-dts";

export default defineConfig({
  plugins: [dts()],
  build: {
    outDir: "./dist",
    copyPublicDir: false,
    lib: {
      // prettier-ignore
      entry: {
        // Note: there is no barrel file to import all two-web/kit packages
        "context": "./context/index.ts",
        "dom-ranges": "./dom-ranges/index.ts",
        "pre-fetcher": "./pre-fetcher/index.ts",
        "route-guards": "./route-guards/index.ts",
        "signals": "./signals/index.ts",
        "spa-router": "./spa-router/index.ts",
        "ssr": "./ssr/index.ts",
        "vite-plugin": "./vite-plugin/index.ts",
      },
      // Since many of the packages are framework agnostic, I want to distribute
      // the compiled JavaScript in as many formats as possible so that it can
      // be consumed by any framework.
      formats: ["es", "cjs"],
    },
  },
});
