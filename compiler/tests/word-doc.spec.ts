import { beforeEach, expect, test } from "vitest";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("word-doc.html")).document;
});

// Word documents are converted to html with pandoc and then compiled like any
// other markup page.
test("should load", () => {
  assertNoErrors(document);
});

test("should render the converted document content", () => {
  const output = document.body.textContent ?? "";

  expect(output).toContain("Heading 1");
  expect(output).toContain("Dot point 1");
  expect(output).toContain("Second item");
  expect(output).toContain("Third item");
});
