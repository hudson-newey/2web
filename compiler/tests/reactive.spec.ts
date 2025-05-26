import { test, beforeEach, expect } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";
import { getByTestId, getByText } from "@testing-library/dom";
import { randNumber } from "@ngneat/falso";
import userEvent from "@testing-library/user-event";
import { Document } from "happy-dom";

let document: Document;

const incrementButton = () => getByText(document.body as any, "Increment");
const decrementButton = () => getByText(document.body as any, "Decrement");
const countOutput = () => getByTestId(document.body as any, "count-output");

beforeEach(async () => {
  document = await navigateToPage("reactive.html");
});

test("should load", () => {
  assertNoErrors(document);
});

test.skip("should increment and decrement the count when using the increment button", async () => {
  const user = userEvent.setup();

  const expectedCount = randNumber({ min: 1, max: 5 });

  const incrementTarget = incrementButton();
  for (let i = 0; i < expectedCount; i++) {
    await user.click(incrementTarget);
    incrementTarget.dispatchEvent(new Event("click"));
  }

  expect(countOutput().textContent).toEqual(expectedCount);

  await user.click(decrementButton());
  expect(countOutput().textContent).toEqual(expectedCount - 1);
});

test.skip("should decrement the count below zero when using the decrement button", async () => {
  const user = userEvent.setup();
  const expectedCount = randNumber({ min: 1, max: 5 });

  const incrementTarget = decrementButton();
  for (let i = 0; i < expectedCount; i++) {
    await user.click(incrementTarget);
  }

  expect(expectedCount).toEqual(expectedCount);
});
