import { readFileSync } from "node:fs";
import { Browser, BrowserErrorCaptureEnum, BrowserFrame } from "happy-dom";

// Compiled pages reference their JavaScript with `src` attributes pointing at
// files in the build output. The test suite doesn't run a web server, so
// happy-dom can't fetch those files (its fetch against http://localhost:5173
// always fails) and none of the compiled runtime would run.
//
// Inlining the compiled scripts into the page lets the compiled runtime
// (reactive event handlers, initial property assignments, etc.) execute inside
// the tests, which makes it possible to test app functionality.
function inlineCompiledScripts(html: string): string {
  return html.replace(
    /<script([^>]*?)\ssrc="([^"]+)"([^>]*)><\/script>/g,
    (match: string, before: string, src: string, after: string) => {
      let content: string;
      try {
        content = readFileSync(`dist/${src}`, "utf-8");
      } catch {
        // If the script file doesn't exist (e.g. a broken build), leave the
        // tag untouched so that the missing file behavior is preserved.
        return match;
      }

      return `<script${before}${after}>${content}</script>`;
    },
  );
}

export async function navigateToPage(
  location: string,
  beforeLoad?: (frame: BrowserFrame) => void,
): Promise<BrowserFrame> {
  const compiledLocation = `dist/${location}`;
  const html = inlineCompiledScripts(readFileSync(compiledLocation, "utf-8"));

  const browser = new Browser({
    settings: {
      errorCapture: BrowserErrorCaptureEnum.processLevel,
      disableJavaScriptEvaluation: false,
    },
  });
  const page = browser.newPage();

  page.url = "http://localhost:5173/" + location;

  // Some pages run code during the initial script evaluation (e.g. an html
  // output that loads its content from a remote function), so tests need a
  // chance to stub the environment (e.g. fetch) before the page runs.
  beforeLoad?.(page.mainFrame);

  page.content = html;

  return page.mainFrame;
}
