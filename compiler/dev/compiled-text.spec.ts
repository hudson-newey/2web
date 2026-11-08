import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { queryByText } from "@testing-library/dom";

beforeEach(async () => {
  await navigateToPage("compiled-text.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("should replace the initial text", () => {
  expect(
    queryByText(
      document as any,
      "This page is set up correctly",
    )
  ).not.toBeNull();
});
