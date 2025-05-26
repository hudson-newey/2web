import { beforeEach, test } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = await navigateToPage("compiled-styles.html");
});

test("should load", () => {
  assertNoErrors(document);
});
