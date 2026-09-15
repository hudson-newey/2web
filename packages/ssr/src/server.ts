import { createSsrRequestHandler } from "./handler";
import { applyServerHardening, loadRouteManifest, mountRpcEndpoints, mountServerRoutes } from "./routes";
import express from "express";
import fs from "node:fs";

// The default directory that contains the compiled client output. The ssr
// server renders the compiled html documents and serves the compiled assets
// from this directory.
const DEFAULT_CLIENT_DIR = "./dist/";

export async function runServer(port: number = Number(process.env.PORT ?? 5173)) {
  const app = express();

  applyServerHardening(app);

  // Compiled server routes (.server.ts files) are mounted when they have been
  // built (see the routes.json manifest in the server output directory).
  if (fs.existsSync("./dist-server/routes.json")) {
    const manifest = loadRouteManifest("./dist-server/");

    await mountServerRoutes(app, "./dist-server/");
    mountRpcEndpoints(app, "./dist-server/", manifest);
  }

  // Mounted without a route path: it must handle every unmatched request.
  // (express 5 rejects the express 4 wildcard route string `"*"`)
  app.use(createSsrRequestHandler(DEFAULT_CLIENT_DIR));

  // Handler errors are reported without leaking stack traces to the client.
  app.use(
    (
      error: unknown,
      _request: express.Request,
      response: express.Response,
      _next: express.NextFunction,
    ) => {
      console.error("[2web] ssr request error:", error);
      response.status(500).end("internal server error");
    },
  );

  app.listen(port);

  console.log(`[2web] ssr server listening on http://localhost:${port}`);
}
