import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

// A @for loop renders its body once per item of a reactive variable array.
// The loop variable (e.g. 'fruit' in '@for (fruit of $fruits)') is bound to
// the current item inside the body, and the whole loop re-renders when the
// source array is reassigned.

let document: Document;

const listItems = () => document.querySelectorAll(".fruit-list li");
const swapButton = () => getByText(document.body as any, "Swap");

beforeEach(async () => {
  document = (await navigateToPage("for-loops.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render the body once per item of the array", () => {
  expect(listItems().length).toBe(3);
  expect(listItems()[0].textContent).toContain("fruit value: apple");
  expect(listItems()[1].textContent).toContain("fruit value: banana");
  expect(listItems()[2].textContent).toContain("fruit value: cherry");
});

test("should re-render the loop when the source array is reassigned", async () => {
  const user = userEvent.setup();

  await user.click(swapButton());

  expect(listItems().length).toBe(2);
  expect(listItems()[0].textContent).toContain("fruit value: x");
  expect(listItems()[1].textContent).toContain("fruit value: y");
});
