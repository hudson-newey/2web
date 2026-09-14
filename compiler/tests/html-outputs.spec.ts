import { beforeEach, expect, test, vi } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { BrowserFrame, Document } from "happy-dom";

// An html output ('[[ expression ]]') renders its expression's value as html,
// where a text output ('{{ expression }}') renders it as text:
//
//   - A reactive expression is bound through an innerHTML property, so the
//     page re-renders when the value changes.
//   - An expression that calls a (remote) function is loaded asynchronously
//     once when the page loads.

let document: Document;
let frame: BrowserFrame;
let fetchMock: ReturnType<typeof vi.fn>;

const variableOutput = () => document.querySelector(".html-var")!;
const callOutput = () => document.querySelector(".html-call")!;
const textOutput = () => document.querySelector(".text-var")!;

beforeEach(async () => {
  // The remote function call fetches its html during the page's initial
  // script evaluation, so the fetch mock has to be installed before the page
  // runs (happy-dom's fetch can't reach the not-running route server).
  fetchMock = vi.fn(async () => ({
    text: async () => "<em>hello from the server</em>",
  }));

  frame = await navigateToPage("html-outputs.html", (frame) => {
    (frame.window as any).fetch = fetchMock;
  });
  document = frame.document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render a reactive variable as html", () => {
  expect(variableOutput().querySelector("strong")).not.toBeNull();
  expect(variableOutput().innerHTML).toContain("from a variable");

  // The same value through a text output is rendered as text (the tags are
  // visible as text instead of being parsed as html).
  expect(textOutput()!.textContent).toContain("<strong>");
});

test("should render a remote function call's html once it resolves", async () => {
  await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));

  const [url, options] = fetchMock.mock.calls[0];

  expect(url).toBe("/_2web/rpc/api/greetings.server.js/greet");
  expect(JSON.parse(options.body)).toEqual(["server"]);

  await vi.waitFor(() =>
    expect(callOutput().querySelector("em")).not.toBeNull(),
  );
  expect(callOutput().innerHTML).toContain("hello from the server");
});

test("should re-render a reactive html variable when it changes", async () => {
  const user = userEvent.setup();

  await user.click(getByText(document.body as any, "Update"));

  expect(variableOutput().querySelector("em")).not.toBeNull();
  expect(variableOutput().innerHTML).toContain("updated html");
});
