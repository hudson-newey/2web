import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { queryByTestId } from "@testing-library/dom";

const toggleableContent = () => queryByTestId(
  document as any,
  "toggleable-content",
);

beforeEach(async () => {
  await navigateToPage("if-always-closed.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("starts closed", () => {
  expect(toggleableContent()).not.toBeDefined();
});
