import { beforeEach, expect, it } from "vitest";
import { getByTestId, getByText } from "@testing-library/dom";
import { navigateToPage } from "./tests/helpers/fixture";
import userEvent from "@testing-library/user-event";
import { assertNoErrors } from "./tests/helpers/assertions";

const outputElement = () => getByTestId(document as any, "output-element");

beforeEach(async () => {
  await navigateToPage("assignment.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("should have the correct initial value", () => {
  expect(outputElement().textContent).toEqual("Hello");
});

it("should change between states", async () => {
  const user = userEvent.setup();

  await user.click(getByText(document as any, "World"));
  expect(outputElement().textContent).toEqual("World");

  await user.click(getByText(document as any, "Me"));
  expect(outputElement().textContent).toEqual("Me");

  await user.click(getByText(document as any, "Everyone!"));
  expect(outputElement().textContent).toEqual("Everyone!");
});
