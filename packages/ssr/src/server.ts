import { createSsrRequestHandler } from "./handler";
import { applyServerHardening, loadRouteManifest, mountRpcEndpoints, mountServerRoutes } from "./routes";
import express from "express";
import fs from "node:fs";
import path from "node:path";

// The default directory that contains the compiled client output. The ssr
// server renders the compiled html documents and serves the compiled assets
// from this directory.
const DEFAULT_CLIENT_DIR = "./dist/";

// The default directory that contains the compiled server routes and the
// route manifest (routes.json).
const DEFAULT_SERVER_DIR = "./dist-server/";

export interface RunServerOptions {
  // The directory the compiled client output is served from. Defaults to the
  // TWO_WEB_CLIENT_DIR environment variable, then "./dist/".
  clientDir?: string;

  // The directory the compiled server routes (routes.json) are mounted from.
  // Defaults to the TWO_WEB_SERVER_DIR environment variable, then
  // "./dist-server/".
  serverDir?: string;
}

export async function runServer(
  port: number = Number(process.env.PORT ?? 5173),
  options: RunServerOptions = {},
) {
  const clientDir = options.clientDir ?? process.env.TWO_WEB_CLIENT_DIR ?? DEFAULT_CLIENT_DIR;
  const serverDir = options.serverDir ?? process.env.TWO_WEB_SERVER_DIR ?? DEFAULT_SERVER_DIR;

  const app = express();

  applyServerHardening(app);

  // Compiled server routes (.server.ts files) are mounted when they have been
  // built (see the routes.json manifest in the server output directory).
  if (fs.existsSync(path.join(serverDir, "routes.json"))) {
    const manifest = loadRouteManifest(serverDir);

    await mountServerRoutes(app, serverDir);
    mountRpcEndpoints(app, serverDir, manifest);
  }

  // Mounted without a route path: it must handle every unmatched request.
  // (express 5 rejects the express 4 wildcard route string `"*"`)
  app.use(createSsrRequestHandler(clientDir));

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
