import { beforeEach, test } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("compiled-property.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});
