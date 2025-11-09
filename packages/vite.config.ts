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
        iron: "./iron/index.ts",
        keyboard: "./keyboard/index.ts",
        "named-strings": "./named-strings/index.ts",
        "pre-fetcher": "./pre-fetcher/index.ts",
        "route-guards": "./route-guards/index.ts",
        signals: "./signals/index.ts",
        ssr: "./ssr/index.ts",
        typescript: "./typescript/index.ts",
        vdom: "./vdom/index.ts",

        // the "vite-plugin" and "typescript" packages are not compiled here
        // because they are intended to be consumed by development environments
        // instead of production environments.
        //
        // "view-transitions" are included as a static copy.
      },
      formats: ["es", "cjs"],
    },
  },
});
