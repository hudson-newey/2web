import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

const countHeading = () => document.querySelector("h1")!;
const incrementButton = () => getByText(document.body as any, "Increment");
const decrementButton = () => getByText(document.body as any, "Decrement");

beforeEach(async () => {
  document = (await navigateToPage("full-html.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

// The page is a full HTML document (it declares its own doctype, html, head,
// and body tags). The compiler must not emit the document structure twice.
test("should emit the document structure exactly once", () => {
  expect(document.querySelectorAll("html").length).toBe(1);
  expect(document.querySelectorAll("body").length).toBe(1);
});

test("should render the initial count", () => {
  expect(countHeading().textContent).toContain("0");
});

test("should increment and decrement the count", async () => {
  const user = userEvent.setup();

  await user.click(incrementButton());
  await user.click(incrementButton());
  expect(countHeading().textContent).toContain("2");

  await user.click(decrementButton());
  expect(countHeading().textContent).toContain("1");
});

// The page imports the footer component; the imported content must be inlined
// into the page (not left as an unprocessed <tag /> selector).
test("should inline imported components", () => {
  expect(document.body.textContent).toContain("2Web Compiler");
});
