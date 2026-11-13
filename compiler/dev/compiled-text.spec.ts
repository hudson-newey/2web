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
  const el = queryByText(document.body, "This page is set up correctly");
  expect(el).not.toBeNull();
});
