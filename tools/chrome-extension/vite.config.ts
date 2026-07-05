import { defineConfig } from "vite";

export default defineConfig({
  appType: "mpa",
  build: {
    outDir: "./dist",
    copyPublicDir: true,
  },
});
