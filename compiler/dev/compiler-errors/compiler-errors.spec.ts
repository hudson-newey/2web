import { beforeEach, it } from "vitest";
import { navigateToPage } from "../tests/helpers/fixture";
import { assertHasErrors } from "../tests/helpers/assertions";

beforeEach(async () => {
  await navigateToPage("compiler-errors.html");
});

it("has errors", () => {
  assertHasErrors();
});
