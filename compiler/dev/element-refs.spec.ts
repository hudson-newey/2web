import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";

const inputElement = () =>
  getByTestId<HTMLInputElement>(document as any, "input-element");
const outputElement = () =>
  getByTestId<HTMLHeadingElement>(document as any, "output-element");

beforeEach(async () => {
  await navigateToPage("element-refs.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("sets the ids of tagged elements", () => {
  expect(inputElement().id).toBe("inputElement");
  expect(outputElement().id).toBe("outputElement");
});

it("has the correct initial output value", () => {
  expect(outputElement().textContent.trim()).toBe("No Content");
});

it("updates the heading when the input changes", async () => {
  const testText = "This is a test 123";
  const user = userEvent.setup();

  await user.type(inputElement(), testText);
  expect(outputElement().textContent.trim(), testText);
});
