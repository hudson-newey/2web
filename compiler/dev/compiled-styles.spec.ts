import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";

beforeEach(async () => {
  await navigateToPage("compiled-styles.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("styles the target element correctly", () => {
  const el = getByTestId(document.body, "target-el");
  expect(el).toHaveStyle({ color: "green" });
});
