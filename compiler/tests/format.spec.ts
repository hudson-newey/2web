import { beforeEach, expect, test } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

let document: Document;

const clicksButton = () => getByText(document.body as any, "Clicks");

beforeEach(async () => {
  document = (await navigateToPage("format.html")).document;
});

// The .2web file format is 2web's own markup format that compiles to html.
test("should load", () => {
  assertNoErrors(document);
});

test("should render the initial click count", () => {
  expect(clicksButton().textContent).toContain("0 Clicks");
});

test("should update the click count when the button is clicked", async () => {
  const user = userEvent.setup();

  await user.click(clicksButton());
  await user.click(clicksButton());

  expect(clicksButton().textContent).toContain("2 Clicks");
});
