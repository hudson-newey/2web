import { Browser, BrowserErrorCaptureEnum } from "happy-dom";
import { expect } from "vitest";

export async function navigateToPage(location: string) {
  const browser = new Browser({
    settings: {
      errorCapture: BrowserErrorCaptureEnum.processLevel,
      enableJavaScriptEvaluation: true,
      suppressInsecureJavaScriptEnvironmentWarning: true,
      enableCSSFileLoading: true,
    },
  });

  const page = browser.newPage();
  await page.goto(`http://localhost:5174/${location}`);
  await page.waitUntilComplete();

  document = page.mainFrame.document as any;

  expect(document.body.innerHTML.length).greaterThan(0);
}
