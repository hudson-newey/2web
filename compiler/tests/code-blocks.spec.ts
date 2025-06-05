import { beforeEach, describe, expect, test } from "vitest";
import { assertNoErrors } from "./helpers/assertions";
import { navigateToPage } from "./helpers/fixture";
import { BrowserWindow, Document } from "happy-dom";
import { getByTestId } from "@testing-library/dom";

let document: Document;
let window: BrowserWindow;

const getElement = (testId: string) =>
  getByTestId(document.body as any, testId);

beforeEach(async () => {
  const browserFrame = await navigateToPage("code-blocks.html");
  document = browserFrame.document;
  window = browserFrame.window;
});

test("should load", () => {
  assertNoErrors(document);
});

test("empty code block", () => {
  expect(getElement("empty-code").textContent).toEqual("");
});

test("code block with empty script", () => {
  expect(getElement("empty-script").textContent).toEqual("<script></script>");
});

test("code block with empty compiled script", () => {
  expect(getElement("empty-compiled-script").textContent).toEqual(
    "<script compiled></script>",
  );
});

test("html code block", () => {
  const expectedText = `<script>
$ count = 0;
</script>

<h1>{{ $count }}</h1>
<h2>{{ $nonExistent }}</h2>
`;

  expect(getElement("html-code").textContent).toEqual(expectedText);
});

test("javascript code block", () => {
  const expectedText = `
$ message = "Hello World!";
`;

  expect(getElement("javascript-code").textContent).toEqual(expectedText);
});

test("compiled script in <pre> block", () => {
  const expectedText = `<script compiled>
$ greeting = "Hello";
</script>

<h1>{{ $greeting }}</h1>
`;

  expect(getElement("compiled-pre").textContent).toEqual(expectedText);
});

test("empty style block", () => {
  expect(getElement("empty-style").textContent).toEqual("<style></style>");
});

test("css in code block should be emitted as text", () => {
  const expectedText = `
  <style>
    .unstyledElement {
      color: red;
    }
  </style>
`;

  expect(getElement("css-code").textContent).toEqual(expectedText);
});

test("preprocessor text nodes", () => {
  // Because the code is not inside a <pre> block, we expect that the text will
  // be collapsed into a single line.
  expect(getElement("text-nodes").textContent).toEqual("<h1>{{ $count }}</h1>");
});
