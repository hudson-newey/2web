import { Browser, BrowserErrorCaptureEnum } from "happy-dom";
import { expect } from "vitest";

export async function navigateToPage(location: string) {
  const browser = new Browser({
    settings: {
      errorCapture: BrowserErrorCaptureEnum.processLevel,
      disableJavaScriptEvaluation: false,
      enableJavaScriptEvaluation: true,
      suppressInsecureJavaScriptEnvironmentWarning: true,
      disableCSSFileLoading: true,
    },
  });

  const page = browser.newPage();
  await page.goto(`http://localhost:4173/${location}`);
  await page.waitUntilComplete();

  document = page.mainFrame.document as any;

  expect(document.body.innerHTML.length).greaterThan(0);
}
