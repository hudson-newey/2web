import { resolve } from "node:path";
import { builtinModules } from "node:module";
import { defineConfig } from "vite";

const nodeExternals = [
  ...builtinModules,
  ...builtinModules.map((m) => `node:${m}`),
];

export default defineConfig({
  resolve: {
    alias: {
      ws: resolve(__dirname, "node_modules/ws/index.js"),
    },
  },
  build: {
    target: "node18",
    outDir: "dist",
    emptyOutDir: true,
    sourcemap: true,
    minify: false,
    lib: {
      entry: resolve(__dirname, "src/index.ts"),
      formats: ["cjs"],
      fileName: () => "extension.js",
    },
    rollupOptions: {
      external: ["vscode", ...nodeExternals],
    },
  },
});
