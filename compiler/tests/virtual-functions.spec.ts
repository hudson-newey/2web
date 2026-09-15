import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

// Compile time virtual functions are imported from the interop types and are
// evaluated by the compiler while it compiles the page:
//
//   - $uid() expands to a unique incrementing id (useful for aria
//     attributes).
//   - $props() expands to the parameters a component instance was passed.
//     Reactive parameters stay reactive.
//   - $env() and $readFile() expand to compile time string constants.

let document: Document;

const priceTag = () => document.querySelector(".price-tag")!;
const bumpButton = () => getByText(document.body as any, "Bump");

beforeEach(async () => {
  document = (await navigateToPage("virtual-functions.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should expand component props", () => {
  expect(priceTag().textContent).toContain("Reactive tag: 3");
});

// The reactive parameter stays reactive: the component re-renders when the
// passed reactive variable changes.
test("should update a component when a reactive prop changes", async () => {
  const user = userEvent.setup();

  await user.click(bumpButton());

  expect(priceTag().textContent).toContain("Reactive tag: 10");
});

test("should expand unique incrementing ids", () => {
  const ids = document.querySelector(".ids")!.textContent!.trim();

  expect(ids).toBe("uid-1 uid-2");

  // The aria attribute got the next id in the document order.
  expect(document.getElementById("uid-demo")!.getAttribute("aria-describedby")).toBe(
    "uid-3",
  );
});

// An environment variable that isn't set expands to an empty string.
test("should expand environment variables at compile time", () => {
  expect(document.querySelector(".env")!.textContent).toBe("");
});

test("should expand compile time file reads", () => {
  expect(document.querySelector(".license")!.textContent).toContain(
    "MIT-style test license",
  );
});
