import type { Request, Response, NextFunction } from "express";
import { Window } from "happy-dom";
import fs from "node:fs";
import path from "node:path";
import type { SsrConfig } from "../config/config";
import sanitizePath from "sanitize-filename";
import { sleep } from "../../../_shared/sleep";
import { mimeTypeForUrl } from "./mimeTypes";

export function createSsrHandler(config: Readonly<SsrConfig>) {
  const currentDirectory = Deno.cwd();
  const servedPath = path.join(currentDirectory, config.path);

  return async (
    req: Readonly<Request>,
    res: Readonly<Response>,
    next: NextFunction,
  ) => {
    const url = req.protocol + '://' + req.get('host') + req.originalUrl;
    // let safePath = path.join(servedPath, sanitizePath(req.path));
    let safePath = path.join(servedPath, req.originalUrl);

    if (safePath.endsWith("/")) {
      safePath = path.join(safePath, "/index.html");
    }

    const mimeType = mimeTypeForUrl(safePath);

    try {
      const template = fs.readFileSync(safePath, "utf-8");

      // If we are not serving a html file, return it without modification
      if (!safePath.endsWith(".html")) {
        res.status(200).set({ "Content-Type": mimeType }).end(template);
        return;
      }

      const window = new Window({ url });
      const document = window.document;

      document.write(template);

      // Waits for async operations such as timers, resource loading and fetch()
      // on the page to complete.
      // Note that this may get stuck when using intervals or a timer in a loop
      // (see IBrowserSettings for ways to mitigate this).
      await Promise.race([window.happyDOM.waitUntilComplete(), sleep(1_000)]);

      const html = window.document.documentElement.outerHTML;
      await window.happyDOM.close();

      res.status(200).set({ "Content-Type": mimeType }).end(html);
    } catch (e) {
      next(e);
    }
  };
}
