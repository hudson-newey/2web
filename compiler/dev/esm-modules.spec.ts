import { beforeEach, test } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";

beforeEach(async () => {
  await navigateToPage("esm-modules.html");
});

test("should have no errors", () => {
  assertNoErrors();
});
