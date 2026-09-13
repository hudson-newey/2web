import { beforeEach, expect, test } from "vitest";
import { navigateToPage } from "./helpers/fixture";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("syntax-errors.html")).document;
});

// The syntax errors fixture contains intentionally broken syntax. The
// compiler must report every construct through the error overlay (with the
// file and position that caused it) instead of panicking or silently
// dropping the page.
test("should render the compiler error overlay", () => {
  const overlay = document.querySelector(".__2_error_overlay");
  expect(overlay).not.toBeNull();
});

test("should render an error entry per broken construct", () => {
  const errorCount = document.getElementsByClassName("__2_error").length;

  // The fixture breaks a reactive variable, an import path, and an @if
  // condition.
  expect(errorCount).toBeGreaterThanOrEqual(3);
});

test("should attribute errors to the file and position", () => {
  const headings = [
    ...document.getElementsByClassName("__2_error_file"),
  ].map((node) => node.textContent ?? "");

  const positioned = headings.filter((heading) => /:\d+:\d+/.test(heading));

  expect(positioned.length).toBeGreaterThanOrEqual(3);
  expect(positioned.every((heading) => heading.includes("syntax-errors.html"))).toBe(true);
});

test("should describe what is wrong with each construct", () => {
  const overlay = document.body.textContent ?? "";

  expect(overlay).toContain("missing an initial value");
  expect(overlay).toContain("import path must be a quoted string");
  expect(overlay).toContain("missing a condition");
});
