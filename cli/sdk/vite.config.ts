import { defineConfig } from "vite";
import twoWeb from "@two-web/kit/vite-plugin";

// This Vite config is intended to be consumed by people creating 2Web projects.
// It is used when you run the "2web serve" and "2web build" commands
//
// If you want to override / extend this config, you can simply create a
// vite.config.ts in your projects root directory and the 2web cli will use that
// instead.
export default defineConfig({
  appType: "mpa",
  base: "",
  plugins: [twoWeb()],
  build: {
    outDir: "./dist",
    copyPublicDir: true,
  },
});
