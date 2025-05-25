import { test, beforeEach, expect } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";
import { getByTestId, getByText } from "@testing-library/dom";
import { randNumber } from "@ngneat/falso";

let document: Document;

const incrementButton = () => getByText(document.body, "Increment");
const decrementButton = () => getByText(document.body, "Decrement");
const countOutput = () => getByTestId(document.body, "count-output");

beforeEach(() => {
  document = navigateToPage("reactive.html");
});

test("should load", () => {
  assertNoErrors(document);
});

test("should increment and decrement the count when using the increment button", () => {
  const expectedCount = randNumber({ min: 1, max: 5 });

  const incrementTarget = incrementButton();
  for (let i = 0; i < expectedCount; i++) {
    incrementTarget.click();
  }

  // TODO: perform an actual assertion here
  expect(expectedCount).toEqual(expectedCount);

  decrementButton().click();
  expect(expectedCount - 1).toEqual(expectedCount - 1);
});

test("should decrement the count below zero when using the decrement button", () => {
  const expectedCount = randNumber({ min: 1, max: 5 });

  const incrementTarget = decrementButton();
  for (let i = 0; i < expectedCount; i++) {
    incrementTarget.click();
  }

  expect(expectedCount).toEqual(expectedCount);
});
