import { defineConfig } from "vite";
import { viteStaticCopy } from "vite-plugin-static-copy";
import dts from "vite-plugin-dts";

export default defineConfig({
  appType: "mpa",
  plugins: [
    dts(),
    viteStaticCopy({
      targets: [
        {
          src: "view-transitions/transitions/",
          dest: "view-transitions/",
        },
      ],
    }),
  ],
  assetsInclude: ["./view-transitions/transitions/fade.css"],
  build: {
    outDir: "./dist",
    copyPublicDir: false,
    rollupOptions: {
      external: ["vite", "happy-dom", "express", "node", /node:/],
    },
    lib: {
      entry: {
        index: "./index.ts",
        animations: "./animations/index.ts",
        "browser-state": "./browser-state/index.ts",
        "event-listener": "./event-listener/index.ts",
        keyboard: "./keyboard/index.ts",
        "pre-fetcher": "./pre-fetcher/index.ts",
        "route-guards": "./route-guards/index.ts",
        signals: "./signals/index.ts",
        ssr: "./ssr/index.ts",
        threads: "./threads/index.ts",
        // "view-transitions": "./view-transitions/index.css",

        // the "vite-plugin" and "typescript" packages are not compiled here
        // because they are intended to be consumed by development environments
        // instead of production environments.
      },
      formats: ["es", "cjs"],
    },
  },
});
