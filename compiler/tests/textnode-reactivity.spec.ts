import { beforeEach, test } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";
import { Document } from "happy-dom";

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("textnode-reactivity.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});
