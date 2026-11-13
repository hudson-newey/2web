import { beforeEach, expect, it } from "vitest";
import { getByTestId, getByText } from "@testing-library/dom";
import { navigateToPage } from "./tests/helpers/fixture";
import { assertNoErrors } from "./tests/helpers/assertions";
import userEvent from "@testing-library/user-event";

const outputElement = () => getByTestId(document.body, "output-element");

beforeEach(async () => {
  await navigateToPage("assignment.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("should have the correct initial value", () => {
  expect(outputElement()).toHaveTextContent("Hello");
});

it("should change between states", async () => {
  const user = userEvent.setup();

  await user.click(getByText(document.body, "World"));
  expect(outputElement()).toHaveTextContent("World");

  await user.click(getByText(document.body, "Me"));
  expect(outputElement()).toHaveTextContent("Me");

  await user.click(getByText(document.body, "Everyone!"));
  expect(outputElement()).toHaveTextContent("Everyone!");
});
