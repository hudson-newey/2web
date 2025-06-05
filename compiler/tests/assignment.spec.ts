import { beforeEach, expect, test } from "vitest";
import { getByTestId, getByText } from "@testing-library/dom";
import { navigateToPage } from "./helpers/fixture";
import userEvent from "@testing-library/user-event";
import { assertNoErrors } from "./helpers/assertions";
import { BrowserFrame, Document } from "happy-dom";

let document: Document;

const worldButton = () => getByText(document as any, "World");
const meButton = () => getByText(document.body as any, "Me");
const everyoneButton = () => getByText(document.body as any, "Everyone!");
const outputElement = () => getByTestId(document.body as any, "output-element");

beforeEach(async () => {
  document = (await navigateToPage("assignment.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should have the correct initial value", () => {
  expect(outputElement().textContent).toEqual("Hello");
});

test("should change between states", async () => {
  const user = userEvent.setup();

  await user.click(worldButton());
  expect(outputElement().textContent).toEqual("World");

  await user.click(meButton());
  expect(outputElement().textContent).toEqual("Me");

  await user.click(everyoneButton());
  expect(outputElement().textContent).toEqual("Everyone");
});
