import { afterAll, beforeAll, describe, expect, test } from "vitest";
import express, { type Express } from "express";
import fs from "node:fs";
import http from "node:http";
import type { AddressInfo } from "node:net";
import path from "node:path";

import {
  applyServerHardening,
  loadRouteManifest,
  mountRpcEndpoints,
  mountServerRoutes,
  RPC_ENDPOINT_PREFIX,
} from "./routes";

// The rpc endpoints expose the functions that a server script exports to the
// compiled client code. A client calls one with a POST request whose body is a
// json array of the function arguments, and the function return value is
// returned as a string.
//
// The endpoints must only be callable for the functions that the compiler
// recorded in the route manifest.

let serverDir: string;
let serverPort = 0;
let app: Express;
let httpServer: http.Server;

const compiledModule = `
export function greet(name) {
  return "Hello, " + name + "!";
}

export function item(id) {
  return { id, label: "item-" + id };
}

export function boom() {
  throw new Error("secret internal failure");
}

export function notExportedForRpc() {
  return "should never be callable";
}

export default (request, response) => response.json({ ok: true });
`;

// A verb route: its default export handles one http method at its route.
const compiledGetRoute = `
export default (request, response) => {
  response.json({ method: request.method, name: request.query.name ?? null });
};
`;

const compiledPostRoute = `
export default (request, response) => {
  response.json({ method: request.method, body: request.body });
};
`;

// Minimal json request helper (node http, so the tests don't depend on the
// fetch implementation of the test environment).
function request(
  method: string,
  url: string,
  body?: string,
): Promise<{ status: number; text: string }> {
  return new Promise((resolve, reject) => {
    const httpRequest = http.request(
      {
        hostname: "127.0.0.1",
        port: serverPort,
        path: url,
        method,
      },
      (response) => {
        let text = "";

        response.setEncoding("utf-8");
        response.on("data", (chunk: string) => (text += chunk));
        response.on("end", () =>
          resolve({ status: response.statusCode ?? 0, text }),
        );
      },
    );

    httpRequest.on("error", reject);

    if (body !== undefined) {
      httpRequest.setHeader("content-type", "application/json");
      httpRequest.write(body);
    }

    httpRequest.end();
  });
}

const rpcUrl = (module: string, functionName: string) =>
  `${RPC_ENDPOINT_PREFIX}/${module}/${functionName}`;

beforeAll(async () => {
  // The temp dir must live inside the package root: the test environment's
  // module loader refuses to import modules from outside of it.
  serverDir = fs.mkdtempSync(path.join(process.cwd(), ".rpc-test-"));

  fs.writeFileSync(path.join(serverDir, "test.js"), compiledModule);
  fs.writeFileSync(path.join(serverDir, "get.js"), compiledGetRoute);
  fs.writeFileSync(path.join(serverDir, "post.js"), compiledPostRoute);
  fs.writeFileSync(
    path.join(serverDir, "routes.json"),
    JSON.stringify({
      routes: [
        { route: "/test", method: "get", file: "get.js" },
        { route: "/submit", method: "post", file: "post.js" },
      ],
      modules: [{ file: "test.js", rpc: ["greet", "item", "boom"] }],
    }),
  );

  app = express();
  applyServerHardening(app);
  await mountServerRoutes(app, serverDir);
  mountRpcEndpoints(app, serverDir, loadRouteManifest(serverDir));

  httpServer = app.listen(0, "127.0.0.1");
  await new Promise<void>((resolve) => httpServer.once("listening", resolve));

  serverPort = (httpServer.address() as AddressInfo).port;
});

afterAll(() => {
  httpServer?.close();
  fs.rmSync(serverDir, { recursive: true, force: true });
});

describe("rpc endpoints", () => {
  test("should return the function return value as a string", async () => {
    const response = await request(
      "POST",
      rpcUrl("test.js", "greet"),
      JSON.stringify(["world"]),
    );

    expect(response.status).toBe(200);
    expect(response.text).toBe("Hello, world!");
  });

  test("should serialize non string return values as json", async () => {
    const response = await request(
      "POST",
      rpcUrl("test.js", "item"),
      JSON.stringify([7]),
    );

    expect(response.status).toBe(200);
    expect(response.text).toBe('{"id":7,"label":"item-7"}');
  });

  test("should only call functions that the manifest records", async () => {
    // `notExportedForRpc` exists on the module but isn't in the manifest.
    const response = await request(
      "POST",
      rpcUrl("test.js", "notExportedForRpc"),
      JSON.stringify([]),
    );

    expect(response.status).toBe(404);
  });

  test("should 404 for unknown modules and functions", async () => {
    const unknownModule = await request(
      "POST",
      rpcUrl("other.js", "greet"),
      JSON.stringify([]),
    );
    const unknownFunction = await request(
      "POST",
      rpcUrl("test.js", "unknown"),
      JSON.stringify([]),
    );

    expect(unknownModule.status).toBe(404);
    expect(unknownFunction.status).toBe(404);
    expect(unknownFunction.text).not.toContain("stack");
  });

  test("should reject non POST requests", async () => {
    const response = await request("GET", rpcUrl("test.js", "greet"));

    expect(response.status).toBe(405);
  });

  test("should reject bodies that aren't json arrays", async () => {
    const response = await request(
      "POST",
      rpcUrl("test.js", "greet"),
      JSON.stringify({ name: "world" }),
    );

    expect(response.status).toBe(400);
  });

  test("should not leak handler failures to the client", async () => {
    const response = await request(
      "POST",
      rpcUrl("test.js", "boom"),
      JSON.stringify([]),
    );

    expect(response.status).toBe(500);
    expect(response.text).toBe('{"error":"rpc call failed"}');
    expect(response.text).not.toContain("secret internal failure");
  });
});

// Verb routes are mounted for the http method their file name declares, and
// are the only http entry points: server modules are never mounted.
describe("verb routes", () => {
  test("should mount the handler for the method the route declares", async () => {
    const response = await request("GET", "/test?name=2web");

    expect(response.status).toBe(200);
    expect(JSON.parse(response.text)).toEqual({
      method: "GET",
      name: "2web",
    });
  });

  test("should not respond to methods the route doesn't declare", async () => {
    const response = await request("POST", "/test");

    expect(response.status).toBe(404);
  });

  test("should pass the parsed request body to post routes", async () => {
    const response = await request(
      "POST",
      "/submit",
      JSON.stringify({ items: [1, 2] }),
    );

    expect(response.status).toBe(200);
    expect(JSON.parse(response.text)).toEqual({
      method: "POST",
      body: { items: [1, 2] },
    });
  });

  test("should not mount server modules as http endpoints", async () => {
    const response = await request("GET", "/test-module");

    // There is no route for the module file; only verb routes are mounted.
    expect(response.status).toBe(404);
  });

  test("should skip routes with an unknown method", async () => {
    const warn = console.warn;
    console.warn = () => {};

    fs.writeFileSync(
      path.join(serverDir, "teapot.js"),
      "export default (request, response) => response.status(418).end();",
    );
    fs.writeFileSync(
      path.join(serverDir, "routes.json"),
      JSON.stringify({
        routes: [{ route: "/tea", method: "brew", file: "teapot.js" }],
        modules: [],
      }),
    );

    const app = express();
    applyServerHardening(app);
    mountServerRoutes(app, serverDir);

    const httpServer = app.listen(0, "127.0.0.1");
    await new Promise<void>((resolve) => httpServer.once("listening", resolve));
    const port = (httpServer.address() as AddressInfo).port;

    const response = await new Promise<{ status: number }>((resolve) => {
      const httpRequest = http.request(
        { hostname: "127.0.0.1", port, path: "/tea", method: "GET" },
        (response) => {
          response.resume();
          response.on("end", () => resolve({ status: response.statusCode ?? 0 }));
        },
      );
      httpRequest.end();
    });

    httpServer.close();

    console.warn = warn;

    // The route was skipped, so express answers with its default 404 (the
    // error middleware turns it into json).
    expect(response.status).toBe(404);
  });
});
