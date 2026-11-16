import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";

const textOutEl = () => getByTestId(document.body, "text-out");
const valueOut1 = () => getByTestId(document.body, "value-out1");
const valueOut2 = () => getByTestId(document.body, "value-out2");
const placeholderOut = () =>
  getByTestId<HTMLInputElement>(document.body, "placeholder-out");

const worldBtn = () => getByTestId(document.body, "world-btn");
const meBtn = () => getByTestId(document.body, "me-btn");
const everyoneBtn = () => getByTestId(document.body, "everyone-btn");

function expectValue(value: string) {
  expect(textOutEl()).toHaveTextContent(`Current message: ${value}`);
  expect(valueOut1()).toHaveValue(value);
  expect(valueOut2()).toHaveValue(value);
  expect(placeholderOut().placeholder).toBe(value);
}

beforeEach(async () => {
  await navigateToPage("multi-assignment.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("has the correct initial value", () => {
  expect(textOutEl()).toHaveTextContent("Current message: World!");
  expect(valueOut1()).toHaveValue("World!");
  expect(valueOut2()).toHaveValue("World!");
  expect(placeholderOut().placeholder).toBe("World!");
});

it("updates correctly", async () => {
  const user = userEvent.setup();

  // Note that pressing the "world" button removes the exclamation point.
  await user.click(worldBtn());
  expectValue("World");

  await user.click(meBtn());
  expectValue("Me");

  // Clicking the same button again has no effect
  await user.click(meBtn());
  expectValue("Me");

  await user.click(everyoneBtn());
  expectValue("Everyone!");

  // We should be allowed to go back to a previous state.
  await user.click(worldBtn());
  expectValue("World");
});
