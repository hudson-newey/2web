import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

// Components can declare named and default slots, and consumers fill them
// with children:
//
//	<Card>
//		<h2 slot="title">My Card</h2>
//		<p>body content</p>
//	</Card>
//
// Named children replace the slot with the same name (the slot attribute is
// removed from the rendered content), remaining children replace the default
// slot, and unfilled slots render their fallback content.

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("slots.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render the named and default slot content", () => {
  const card = getByText(document.body as any, "body content").closest(".card")!;
  const header = card.previousElementSibling!;

  expect(header.textContent).toContain("My Card");
  expect(card.textContent).toContain("body content");
});

test("should remove the slot attribute from the rendered content", () => {
  expect(document.querySelector("[slot]")).toBeNull();
  expect(document.querySelector("slot")).toBeNull();
});

test("should render the fallback content of unfilled slots", () => {
  const fallbackCard = getByText(document.body as any, "fallback body").closest(".card")!;
  const fallbackHeader = fallbackCard.previousElementSibling!;

  expect(fallbackHeader.textContent).toContain("fallback title");
});

// The slot content is the consumer's own markup, so its reactive variables
// keep working inside the component.
test("should keep reactive variables working inside slot content", async () => {
  const user = userEvent.setup();

  expect(getByText(document.body as any, "0")).toBeTruthy();

  await user.click(getByText(document.body as any, "Bump"));

  expect(getByText(document.body as any, "42")).toBeTruthy();
});
