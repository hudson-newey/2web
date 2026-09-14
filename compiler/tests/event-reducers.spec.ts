import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

// Regression tests for event reducer bugs:
//
//   - A cross variable assignment ('@click="$x = $y + 1"') used to emit a
//     broken listener (an empty runtime variable name for the read variable,
//     a second listener wired by the read variable, and the raw selector
//     instead of its runtime value).
//   - A "=" character inside a reducer's value expression
//     ('@click="$text = \'a=b\'"') used to truncate the expression (the
//     reducer was split at every "=" instead of the assignment operator).
//   - A comparison on the right hand side ('@click="$flag = $x == 1"') used
//     to lose the comparison entirely.
//   - Compound assignments ('@click="$count += 2"') were not supported.

let document: Document;

const paragraph = (label: string) =>
  getByText(document.body as any, label).closest("p")!;
const click = async (label: string) => {
  const user = userEvent.setup();
  await user.click(getByText(document.body as any, label));
};

beforeEach(async () => {
  document = (await navigateToPage("event-reducers.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render the initial values", () => {
  expect(paragraph("X").textContent).toContain("X 0");
  expect(paragraph("Y").textContent).toContain("Y 10");
  expect(paragraph("Count").textContent).toContain("Count 5");
});

// $x is assigned from $y's runtime value through $x's update cascade.
test("should assign across variables", async () => {
  await click("Cross");

  expect(paragraph("X").textContent).toContain("X 11");
  // The read variable's value must not change.
  expect(paragraph("Y").textContent).toContain("Y 10");
});

test("should keep the value expression intact when it contains a =", async () => {
  await click("Equals");

  expect(paragraph("Text").textContent).toContain("Text a=b");
});

test("should evaluate comparisons on the right hand side", async () => {
  // $x is 11 after the cross variable assignment, so the comparison
  // evaluates to true (the comparison itself used to be dropped entirely).
  await click("Cross");
  expect(paragraph("Flag").textContent).toContain("Flag false");

  await click("Compare");

  expect(paragraph("Flag").textContent).toContain("Flag true");
});

test("should support compound assignments", async () => {
  await click("Add");

  expect(paragraph("Count").textContent).toContain("Count 7");
});

// A static variable's text output is inlined into the page at compile time
// (escaped, matching the text rendering semantics of a runtime text output).
test("should inline a static text output into the page", () => {
  const staticOutput = () => document.querySelector(".static-out")!;

  expect(staticOutput().innerHTML).toContain("rendered at compile time");
  expect(staticOutput().querySelector("span")).toBeNull();
});
