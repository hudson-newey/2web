# 2Web Vite Plugin

A [Vite](https://vite.dev/) plugin to make 2Web development painless.

```ts
import { defineConfig } from "vite";
import twoWeb from "@two-web/kit/vite-plugin";

export default defineConfig({
  appType: "mpa",
  plugins: [twoWeb()],
});
```
