import { createServer as createViteServer } from "vite";
import { handleSsrRequest } from "./handler";
import { applyServerHardening, mountServerRoutes } from "./routes";
import express from "express";
import fs from "node:fs";

export async function runServer(port: number = 5173) {
  const app = express();

  applyServerHardening(app);

  // Compiled server routes (.server.ts files) are mounted when they have been
  // built (see the routes.json manifest in the server output directory).
  if (fs.existsSync("./dist-server/routes.json")) {
    await mountServerRoutes(app, "./dist-server/");
  }

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

  app.use("*", handleSsrRequest);

  app.listen(port);
}
