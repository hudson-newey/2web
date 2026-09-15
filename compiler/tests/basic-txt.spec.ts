import { beforeEach, expect, test } from "vitest";
import { navigateToPage } from "./helpers/fixture";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("basic-txt.txt")).document;
});

// Text files are passed through to the build output without compilation. The
// full source (including the text that looks like markup) must be present in
// the output.
test("should pass the text file through unmodified", () => {
  const output = document.body.textContent ?? "";

  expect(output).toContain("This is a test");
  expect(output).toContain("html elements should not be rendered");
  expect(output).toContain("TODO: We should be able to import this file");
});
