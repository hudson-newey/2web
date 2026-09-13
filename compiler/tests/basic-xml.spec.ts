import { beforeEach, expect, test } from "vitest";
import { navigateToPage } from "./helpers/fixture";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("basic-xml.xml")).document;
});

// Browsers can render xml natively, so xml files are passed through to the
// build output without compilation.
test("should pass the xml file through with its content", () => {
  const output = document.body.textContent ?? "";

  expect(output).toContain("Learning XML");
  expect(output).toContain("Erik T. Ray");
  expect(output).toContain("XML in a Nutshell");
  expect(output).toContain("Beginning XML");
});
