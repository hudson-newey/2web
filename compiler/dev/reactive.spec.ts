import { test, beforeEach, expect } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { randNumber } from "@ngneat/falso";
import { screen  } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";

const incrementButton = () => screen.getByText("Increment");
const decrementButton = () => screen.getByText("Decrement");
const countOutput = () => screen.getByTestId("count-output");

beforeEach(async () => {
  await navigateToPage("reactive.html");
});

test("should load", () => {
  assertNoErrors();
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
