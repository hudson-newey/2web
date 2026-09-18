import type { Request, Response, NextFunction } from "express";
import { fileURLToPath } from "node:url";
import { Window } from "happy-dom";
import fs from "node:fs";
import path from "node:path";
import type { SsrConfig } from "../config/config";
import sanitizePath from "sanitize-filename";

export function createSsrHandler(config: Readonly<SsrConfig>) {
  const currentDirectory = path.dirname(fileURLToPath(import.meta.url));
  const servedPath = path.join(currentDirectory, config.path);

  return async (
    req: Readonly<Request>,
    res: Readonly<Response>,
    next: NextFunction,
  ) => {
    const url = req.originalUrl;
    const safePath = sanitizePath(path.resolve(servedPath + req.path));

    try {
      const template = fs.readFileSync(safePath, "utf-8");

      // If we are not serving a html file, return it without modification
      if (!req.path.endsWith(".html") && !req.path.endsWith("/")) {
        res.status(200).end(template);
        return;
      }

      const window = new Window({ url });
      const document = window.document;

      document.write(template);

      // Waits for async operations such as timers, resource loading and fetch() on the page to complete
      // Note that this may get stuck when using intervals or a timer in a loop (see IBrowserSettings for ways to mitigate this)
      await window.happyDOM.waitUntilComplete();
      const html = window.document.documentElement.outerHTML;

      res.status(200).set({ "Content-Type": "text/html" }).end(html);

      await window.happyDOM.close();
    } catch (e) {
      next(e);
    }
  };
}
