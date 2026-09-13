import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

const conditionOutput = () => getByText(document.body as any, "Hello World!");

beforeEach(async () => {
  document = (await navigateToPage("control-flow-if.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

// The @if block should compile into a span that is wired to the $isOpen
// reactive variable through a runtime selector. No compile time selector
// (e.g. *hidden="$isOpen") may be left in the page.
test("should compile the @if block into a wired element", () => {
  const conditionalElement = conditionOutput().closest("span")!;

  expect(conditionalElement).not.toBeNull();
  // The element id depends on the build wide selector allocation order, so
  // match on the selector prefix instead of an exact id.
  expect(conditionalElement.outerHTML).toMatch(/data-__2_element_\d+/);
  expect(conditionalElement.innerHTML).not.toContain("*hidden");
});

// $isOpen is declared as false, so the @if block should be hidden on the
// initial render.
test("should hide the conditional content when the condition is false", () => {
  const conditionalElement = conditionOutput().closest("span")!;

  expect(conditionalElement.hidden).toBe(true);
});
