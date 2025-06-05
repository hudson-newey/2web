import { readFileSync } from "node:fs";
import { Browser, BrowserErrorCaptureEnum, BrowserFrame } from "happy-dom";

export async function navigateToPage(location: string): Promise<BrowserFrame> {
  const compiledLocation = `dist/${location}`;
  const html = readFileSync(compiledLocation, "utf-8");

  const browser = new Browser({
    settings: {
      errorCapture: BrowserErrorCaptureEnum.processLevel,
      disableJavaScriptEvaluation: false,
    },
  });
  const page = browser.newPage();

  page.url = "http://localhost:5173/" + location;
  page.content = html;

  return page.mainFrame;
}
