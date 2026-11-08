import { beforeEach, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";

beforeEach(async () => {
  await navigateToPage("textnode-reactivity.html");
});

it("has no errors", () => {
  assertNoErrors();
});
