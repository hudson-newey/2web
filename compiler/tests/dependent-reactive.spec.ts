import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

const increaseButton = () => getByText(document.body as any, "Increase");
const countParagraph = () => getByText(document.body as any, "Count").closest("p")!;
const doubledParagraph = () => getByText(document.body as any, "Doubled").closest("p")!;
const tripledParagraph = () => getByText(document.body as any, "Tripled plus one").closest("p")!;

beforeEach(async () => {
  document = (await navigateToPage("dependent-reactive.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render the initial values of computed variables", () => {
  // $count = 1, $doubled = $count * 2, $tripledPlusOne = $doubled * 3 + 1
  expect(countParagraph().textContent).toContain("Count 1");
  expect(doubledParagraph().textContent).toContain("Doubled 2");
  expect(tripledParagraph().textContent).toContain("Tripled plus one 7");
});

// Changing $count must re-evaluate every variable computed from it (through
// the whole dependency chain) and update their bound elements.
test("should recompute dependent variables when the source variable changes", async () => {
  const user = userEvent.setup();

  await user.click(increaseButton());
  expect(countParagraph().textContent).toContain("Count 2");
  expect(doubledParagraph().textContent).toContain("Doubled 4");
  expect(tripledParagraph().textContent).toContain("Tripled plus one 13");

  await user.click(increaseButton());
  expect(countParagraph().textContent).toContain("Count 3");
  expect(doubledParagraph().textContent).toContain("Doubled 6");
  expect(tripledParagraph().textContent).toContain("Tripled plus one 19");
});
