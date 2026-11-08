import { beforeEach, test } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { Document } from "happy-dom";

beforeEach(async () => {
  await navigateToPage("control-flow-if.html");
});

test("should have no errors", () => {
  assertNoErrors();
});
