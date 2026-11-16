import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";

const inputEl = () => getByTestId(document.body, "input-el");
const outputEl = () => getByTestId(document.body, "output-el");

beforeEach(async () => {
  await navigateToPage("input-reflection.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("has the correct initial value", () => {
  expect(outputEl()).toHaveTextContent("");
});

it("updates when the input changes", async () => {
  const user = userEvent.setup({ document });

  await user.type(inputEl(), "This is a test");
  expect(outputEl()).toHaveTextContent("This is a test");

  // Test clearing the input
  await user.clear(inputEl());
  expect(outputEl()).toHaveTextContent("");
});
