import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";

beforeEach(async () => {
  await navigateToPage("compiled-property.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("input should have the correct initial value", () => {
  const inputElement = getByTestId<HTMLInputElement>(document as any, "value-output");
  expect(inputElement.value).toBe("Hello World!");
});
