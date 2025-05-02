import { readFileSync } from "node:fs";
import { JSDOM } from "jsdom";

export function navigateToPage(location: string): Document {
  const compiledLocation = `dist/${location}`;
  const html = readFileSync(compiledLocation, "utf-8");
  const dom = new JSDOM(html);

  const document = dom.window.document;

  return document;
}
