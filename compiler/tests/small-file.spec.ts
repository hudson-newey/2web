import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

const incrementButton = () => getByText(document.body as any, "Click Me");
const toggleButton = () => getByText(document.body as any, "Toggle");
const greetingContainer = () =>
  getByText(document.body as any, "Hello").closest("div")!;

beforeEach(async () => {
  document = (await navigateToPage("small-file.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render the initial greeting count", () => {
  expect(greetingContainer().textContent).toContain("Hello 0");
});

// The counter is a self-referencing assignment ($greeting = $greeting + 1),
// which requires the compiled runtime to keep track of state between events.
test("should increment the greeting count when the button is clicked", async () => {
  const user = userEvent.setup();

  await user.click(incrementButton());
  await user.click(incrementButton());
  await user.click(incrementButton());

  expect(greetingContainer().textContent).toContain("Hello 3");
});

// The toggle drives an @if block through a negated assignment
// ($toggled = !$toggled), so the dropdown content should appear and disappear.
test("should toggle the conditional dropdown when the toggle button is clicked", async () => {
  const user = userEvent.setup();

  const dropdown = getByText(document.body as any, "This is a dropdown!").closest("span")!;
  expect(dropdown.hidden).toBe(true);

  await user.click(toggleButton());
  expect(dropdown.hidden).toBe(false);

  await user.click(toggleButton());
  expect(dropdown.hidden).toBe(true);
});
