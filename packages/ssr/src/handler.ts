import type { Request, Response, NextFunction } from "express";
import { Window } from "happy-dom";
import fs from "node:fs";
import path from "node:path";

// The directory index that directory requests (e.g. "/") resolve to.
const DIRECTORY_INDEX = "index.html";

// Creates the express request handler that serves the compiled client output:
//
//   - html documents are server side rendered with happy-dom (the page is
//     parsed and executed, and the resulting dom is serialized back), and
//   - every other file is served statically from the compiled client output
//     directory (with the content type express infers from the file path).
//
// The clientDir is the directory the compiler emitted the client build into
// (e.g. "./dist/").
export function createSsrRequestHandler(clientDir: string) {
  return async (req: Request, res: Response, next: NextFunction) => {
    try {
      const requestPath = decodeURIComponent(req.path);
      const resolvedPath = path.join(clientDir, requestPath);

      if (isHtmlRequest(requestPath)) {
        await renderHtmlDocument(resolvedPath, res);
        return;
      }

      serveStaticFile(resolvedPath, res);
    } catch (error) {
      next(error);
    }
  };
}

function isHtmlRequest(requestPath: string): boolean {
  return (
    requestPath.endsWith("/") ||
    requestPath === "" ||
    requestPath.endsWith(".html")
  );
}

// Server renders one compiled html document. Directory requests (e.g. "/")
// resolve to the directory index document.
async function renderHtmlDocument(resolvedPath: string, res: Response) {
  const documentPath = resolvedPath.endsWith("/")
    ? path.join(resolvedPath, DIRECTORY_INDEX)
    : resolvedPath;

  if (!fs.existsSync(documentPath)) {
    res.status(404).set({ "Content-Type": "text/html" }).end(notFoundPage());
    return;
  }

  const template = fs.readFileSync(documentPath, "utf-8");

  const window = new Window({ url: "https://localhost/" });
  const document = window.document;

  document.write(template);

  // Waits for async operations such as timers, resource loading and fetch() on the page to complete
  // Note that this may get stuck when using intervals or a timer in a loop (see IBrowserSettings for ways to mitigate this)
  try {
    await window.happyDOM.waitUntilComplete();
    const html = window.document.documentElement.outerHTML;

    res.status(200).set({ "Content-Type": "text/html" }).end(html);
  } finally {
    await window.happyDOM.close();
  }
}

// Serves one compiled client asset. Missing files respond with a 404 instead
// of the express error page.
function serveStaticFile(resolvedPath: string, res: Response) {
  if (!fs.existsSync(resolvedPath) || !fs.statSync(resolvedPath).isFile()) {
    res.status(404).json({ error: "not found" });
    return;
  }

  // sendFile sets the content type from the file extension.
  res.sendFile(path.resolve(resolvedPath));
}

function notFoundPage(): string {
  return "<!doctype html><html><body><h1>404 Not Found</h1></body></html>";
}
