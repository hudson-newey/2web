import { test } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";

test("should load", () => {
  const document = navigateToPage("compiled-styles.html");
  assertNoErrors(document);
});
