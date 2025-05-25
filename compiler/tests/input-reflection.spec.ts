import { beforeEach, test } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";

let document: Document;

beforeEach(() => {
  document = navigateToPage("input-reflection.html");
});

test("should load", () => {
  assertNoErrors(document);
});
