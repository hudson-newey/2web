import express, { type Express } from "express";
import fs from "node:fs";
import http from "node:http";
import path from "node:path";
import { pathToFileURL } from "node:url";

export interface ServerRoute {
  route: string;
  file: string;

  // The exported functions that the generated rpc endpoints call.
  rpc?: string[];
}

export interface ServerRouteManifest {
  routes: ServerRoute[];
}

export interface RouteServerOptions {
  // The directory that contains the compiled server routes and the route
  // manifest (routes.json). Defaults to "./dist-server/".
  serverDir?: string;

  // The directory that contains the compiled client output. When set, the
  // compiled client output is served statically.
  clientDir?: string;

  port?: number;
}

// The default express request timeout, header read timeout, and keep alive
// timeout are generous; tightening them protects the server from slow loris
// style attacks without breaking legitimate slow clients.
const SERVER_TIMEOUTS = {
  requestTimeout: 30_000,
  headersTimeout: 20_000,
  keepAliveTimeout: 5_000,
  connectionsCheckingInterval: 1_000,
};

// The maximum size of a request body. Server routes that need larger bodies
// should parse the body themselves.
const MAX_BODY_SIZE = "100kb";

export function loadRouteManifest(serverDir: string): ServerRouteManifest {
  const manifestPath = path.join(serverDir, "routes.json");

  if (!fs.existsSync(manifestPath)) {
    return { routes: [] };
  }

  try {
    return JSON.parse(fs.readFileSync(manifestPath, "utf-8"));
  } catch {
    // A corrupt manifest degrades to no routes; the compiler regenerates it
    // on the next build.
    return { routes: [] };
  }
}

// Applies the default security hardening to the generated express server.
//
// The defaults are intentionally dependency free (no helmet required) and
// cover the header and parser configuration that a publicly exposed express
// server should never run without:
//
//   - The "x-powered-by" header is disabled so the server doesn't advertise
//     its technology stack.
//   - MIME sniffing is disabled so browsers can't be tricked into
//     interpreting a response as executable.
//   - Clickjacking is mitigated with a SAMEORIGIN frame policy.
//   - Referrers are trimmed down to the origin.
//   - Strict transport security is only advertised over secure connections.
//   - Request bodies are size limited to protect against denial of service.
export function applyServerHardening(app: Express): void {
  app.disable("x-powered-by");

  app.use((_request, response, next) => {
    response.setHeader("X-Content-Type-Options", "nosniff");
    response.setHeader("X-Frame-Options", "SAMEORIGIN");
    response.setHeader("Referrer-Policy", "strict-origin-when-cross-origin");

    if (response.locals?.secure ?? _request.secure) {
      response.setHeader(
        "Strict-Transport-Security",
        "max-age=31536000; includeSubDomains",
      );
    }

    next();
  });

  app.use(express.json({ limit: MAX_BODY_SIZE }));
  app.use(express.urlencoded({ extended: false, limit: MAX_BODY_SIZE }));
}

// Mounts every compiled server route on the given express app.
//
// Each compiled handler module default exports an express request handler.
// The handlers are imported once at startup; a handler that crashes only
// breaks its own route (express catches handler errors and forwards them to
// the error middleware).
export async function mountServerRoutes(
  app: Express,
  serverDir: string,
): Promise<void> {
  const manifest = loadRouteManifest(serverDir);

  for (const serverRoute of manifest.routes) {
    const modulePath = path.join(serverDir, serverRoute.file);
    const module = await import(pathToFileURL(modulePath).href);

    const handler = module.default;

    if (typeof handler !== "function") {
      console.warn(
        `[2web] server route '${serverRoute.route}' (${serverRoute.file}) does not default export a function and was skipped`,
      );
      continue;
    }

    app.all(serverRoute.route, handler);
  }
}

// Starts the generated express server with the compiled server routes.
//
// Returns the node http server so that callers (e.g. the cli) can stop it
// programmatically.
export async function startRouteServer(
  options: RouteServerOptions = {},
): Promise<http.Server> {
  const serverDir = options.serverDir ?? "./dist-server/";
  const port = options.port ?? Number(process.env.PORT ?? 3000);

  const app = express();
  applyServerHardening(app);

  const manifest = loadRouteManifest(serverDir);

  await mountServerRoutes(app, serverDir);
  mountRpcEndpoints(app, serverDir, manifest);

  if (options.clientDir) {
    // Serve the compiled client output statically. The index option is
    // disabled so that static serving can't shadow the server routes.
    app.use(
      express.static(options.clientDir, {
        index: false,
        maxAge: "1h",
      }),
    );
  }

  const server = http.createServer(
    {
      requestTimeout: SERVER_TIMEOUTS.requestTimeout,
      headersTimeout: SERVER_TIMEOUTS.headersTimeout,
      keepAliveTimeout: SERVER_TIMEOUTS.keepAliveTimeout,
      connectionsCheckingInterval: SERVER_TIMEOUTS.connectionsCheckingInterval,
    },
    app,
  );

  // Requests that don't match a server route or a client file are rejected
  // with a json 404 (the generated server is an api server by default).
  app.use((_request, response) => {
    response.status(404).json({ error: "not found" });
  });

  // Handler errors are reported without leaking stack traces to the client.
  app.use(
    (
      error: unknown,
      _request: express.Request,
      response: express.Response,
      _next: express.NextFunction,
    ) => {
      console.error("[2web] server route error:", error);
      response.status(500).json({ error: "internal server error" });
    },
  );

  server.listen(port);

  await new Promise<void>((resolve) => {
    server.once("listening", resolve);
  });

  console.log(`[2web] server listening on http://localhost:${port}`);

  // Graceful shutdown on the standard process signals so that the cli (and
  // process managers) can stop the server without dropping in flight
  // requests.
  for (const signal of ["SIGTERM", "SIGINT"] as const) {
    process.on(signal, () => {
      server.close(() => process.exit(0));
    });
  }

  return server;
}

// ---------------------------------------------------------------------------
// Rpc endpoints
// ---------------------------------------------------------------------------

// The url prefix that the generated rpc endpoints are mounted under.
// e.g. "/_2web/rpc/api/users.server.js/greet"
export const RPC_ENDPOINT_PREFIX = "/_2web/rpc";

// Mounts the generated rpc endpoints on the given express app.
//
// A client calls an exported server function through a POST request whose body
// is a json array of the function arguments:
//
//	POST /_2web/rpc/api/users.server.js/greet   body: ["world"]
//
// The endpoint dispatches to the compiled function and returns the function
// return value as a string.
//
// Only functions that the compiler recorded in the route manifest are
// callable: the dispatcher validates the module and function against the
// manifest, so a request can never execute something that wasn't exported at
// compile time.
export function mountRpcEndpoints(
  app: Express,
  serverDir: string,
  manifest: ServerRouteManifest,
): void {
  app.use(RPC_ENDPOINT_PREFIX, (request, response, next) => {
    void handleRpcRequest(request, response, next, serverDir, manifest);
  });
}

async function handleRpcRequest(
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
  serverDir: string,
  manifest: ServerRouteManifest,
): Promise<void> {
  if (request.method !== "POST") {
    response.status(405).json({ error: "rpc endpoints only accept POST" });
    return;
  }

  // The request path is "/<module path>/<function name>".
  const trimmed = request.path.replace(/^\//, "");
  const functionSeparator = trimmed.lastIndexOf("/");

  if (functionSeparator <= 0) {
    response.status(404).json({ error: "unknown rpc endpoint" });
    return;
  }

  const modulePath = trimmed.slice(0, functionSeparator);
  const functionName = trimmed.slice(functionSeparator + 1);

  // Only functions that the compiler recorded for this module can be called.
  const route = manifest.routes.find((entry) => entry.file === modulePath);

  if (!route || !(route.rpc ?? []).includes(functionName)) {
    response.status(404).json({ error: "unknown rpc endpoint" });
    return;
  }

  const args = request.body;

  if (!Array.isArray(args)) {
    response.status(400).json({ error: "rpc arguments must be a json array" });
    return;
  }

  try {
    const modulePathUrl = pathToFileURL(
      path.join(serverDir, modulePath),
    ).href;

    const module = await import(/* @vite-ignore */ modulePathUrl);
    const target = module[functionName];

    if (typeof target !== "function") {
      response.status(404).json({ error: "unknown rpc endpoint" });
      return;
    }

    // The function return is returned as a string. Strings are passed through
    // verbatim; every other value is serialized as json so that objects and
    // arrays don't degrade to "[object Object]".
    const result = await target(...args);
    const body =
      typeof result === "string" ? result : (JSON.stringify(result) ?? "");

    response.status(200).send(body);
  } catch (error) {
    console.error(`[2web] rpc call to '${modulePath}/${functionName}' failed:`, error);
    response.status(500).json({ error: "rpc call failed" });
  }
}
