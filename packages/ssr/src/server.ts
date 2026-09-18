import { createServer as createViteServer } from "vite";
import { createSsrHandler } from "./renderer/handler";
import express from "express";
import type { SsrConfig } from "./config/config";
import { mergeDefaultConfig } from "./config/defaultConfig";

export async function runServer(userConfig: Readonly<Partial<SsrConfig>>) {
  const config = mergeDefaultConfig(userConfig);

  const app = express();

  // Create Vite server in middleware mode and configure the app type as
  // 'custom', disabling Vite's own HTML serving logic so parent server
  // can take control
  const vite = await createViteServer({
    server: { middlewareMode: true },
    appType: "custom",
  });

  // Use vite's connect instance as middleware. If you use your own
  // express router (express.Router()), you should use router.use
  // When the server restarts (for example after the user modifies
  // vite.config.js), `vite.middlewares` is still going to be the same
  // reference (with a new internal stack of Vite and plugin-injected
  // middlewares). The following is valid even after restarts.
  app.use(vite.middlewares);

  const handler = createSsrHandler(config);

  app.use(/(.*)/, handler);

  app.listen(config.port);
}
