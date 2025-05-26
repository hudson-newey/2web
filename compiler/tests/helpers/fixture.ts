import { readFileSync } from "node:fs";
import { Browser, BrowserErrorCaptureEnum } from "happy-dom";

export async function navigateToPage(location: string) {
  const compiledLocation = `dist/${location}`;
  const html = readFileSync(compiledLocation, "utf-8");

  const browser = new Browser({
    settings: {
      errorCapture: BrowserErrorCaptureEnum.processLevel,
      disableJavaScriptEvaluation: false,
    },
  });
  const page = browser.newPage();

  // Navigates page
  await page.goto("http://localhost:5173/" + location);

  const document = page.mainFrame.document;

  return document;
}
