import { Browser, BrowserErrorCaptureEnum, BrowserFrame } from "happy-dom";

export async function navigateToPage(location: string): Promise<BrowserFrame> {
  const browser = new Browser({
    settings: {
      errorCapture: BrowserErrorCaptureEnum.processLevel,
      disableJavaScriptEvaluation: false,
    },
  });

  const page = browser.newPage();
  await page.goto(location);
  await page.waitUntilComplete();

  return page.mainFrame;
}
