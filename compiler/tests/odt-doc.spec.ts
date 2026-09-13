import { beforeEach, expect, test } from "vitest";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("odt-doc.html")).document;
});

// OpenDocument text files are converted to html with pandoc and then compiled
// like any other markup page.
test("should load", () => {
  assertNoErrors(document);
});

test("should render the converted document content", () => {
  const output = document.body.textContent ?? "";

  expect(output).toContain("Hello World!");
  expect(output).toContain("Supporting .odt files");
  expect(output).toContain("second dot point");
});
