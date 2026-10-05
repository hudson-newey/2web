import { fileURLToPath } from "node:url";
import { Window } from "happy-dom";
import type { SsrConfig } from "../config/config.ts";
import { join, dirname, resolve } from "@std/path";
import sanitizePath from "sanitize-filename";

export function createSsrHandler(config: Readonly<SsrConfig>) {
  const currentDirectory = dirname(fileURLToPath(import.meta.url));
  const servedPath = join(currentDirectory, config.path);

  return async (req: Readonly<Request>) => {
    const url = new URL(req.url);
    const safePath = sanitizePath(resolve(servedPath + url.pathname));

    try {
      const template = await Deno.readTextFile(safePath);

      // If we are not serving a html file, return it without modification
      const href = url.href;
      if (!href.endsWith(".html") && !href.endsWith("/")) {
        return new Response(template, {
          status: 200,
        });
      }

      const window = new Window({ url: href });

      const document = window.document;
      document.write(template);

      // Waits for async operations such as timers, resource loading and fetch() on the page to complete
      // Note that this may get stuck when using intervals or a timer in a loop (see IBrowserSettings for ways to mitigate this)
      await window.happyDOM.waitUntilComplete();
      const html = window.document.documentElement.outerHTML;

      await window.happyDOM.close();

      return new Response(html, {
        status: 200,
        headers: {
          "Content-Type": "text/html",
        },
      });
    } catch (e) {
      console.error(e);
    }

    return new Response("Unknown Error", {
      status: 500,
    });
  };
}
