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
  //
  // The mock simulates the server route: the greeting echoes its argument,
  // so that re-issued calls render a distinct value.
  fetchMock = vi.fn(async (_url: unknown, options?: { body: string }) => {
    const name = JSON.parse(options?.body ?? "[]")[0] ?? "";

    return { text: async () => "Hello, " + name + "!" };
  });

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
  // (the mock echoes the call's argument, see beforeEach)
  // The page also loads a computed variable that calls the same function, so
  // the call for this output is matched by its (expanded) argument.
  await vi.waitFor(() =>
    expect(
      fetchMock.mock.calls.some((call) => JSON.parse(call[1].body)[0] === "server"),
    ).toBe(true),
  );

  const [url, options] = fetchMock.mock.calls.find(
    (call) => JSON.parse(call[1].body)[0] === "server",
  )!;

  expect(url).toBe("/_2web/rpc/api/greetings.server.js/greet");

  // The reactive variable argument is expanded to its current value.
  expect(JSON.parse(options.body)).toEqual(["server"]);

  await vi.waitFor(() =>
    expect(callOutput().innerHTML).toContain("Hello, server!"),
  );
});

// The load function of an html output that calls a remote function with a
// reactive variable is re-run by the variable's update cascade, so the call
// is re-issued (with the new value) whenever the variable changes.
test("should re-call the remote function when a reactive argument changes", async () => {
  const user = userEvent.setup();

  await vi.waitFor(() => expect(fetchMock.mock.calls.length).toBeGreaterThan(0));

  await user.click(getByText(document.body as any, "Rename"));

  // The output's load function is re-run by $name's update cascade, so the
  // call is re-issued with the new value.
  await vi.waitFor(() =>
    expect(
      fetchMock.mock.calls.some((call) => JSON.parse(call[1].body)[0] === "other"),
    ).toBe(true),
  );

  await vi.waitFor(() =>
    expect(callOutput().innerHTML).toContain("Hello, other!"),
  );
});

// A computed variable whose value expression calls a remote function is
// evaluated asynchronously (its runtime value is a placeholder until the
// call resolves) and is re-evaluated whenever a variable it is computed from
// changes.
test("should evaluate a computed variable that calls a remote function", async () => {
  const computedOutput = () => document.querySelector(".computed-call")!;

  await vi.waitFor(() =>
    expect(computedOutput().textContent).toContain("Hello, 0!"),
  );

  const user = userEvent.setup();

  const callCountBefore = fetchMock.mock.calls.length;

  await user.click(getByText(document.body as any, "Count"));

  // The call is re-issued with the new value of $count.
  await vi.waitFor(() =>
    expect(fetchMock.mock.calls.length).toBeGreaterThan(callCountBefore),
  );
  expect(
    JSON.parse(fetchMock.mock.calls[fetchMock.mock.calls.length - 1][1].body),
  ).toEqual([1]);

  await vi.waitFor(() =>
    expect(computedOutput().textContent).toContain("Hello, 1!"),
  );
});

test("should re-render a reactive html variable when it changes", async () => {
  const user = userEvent.setup();

  await user.click(getByText(document.body as any, "Update"));

  expect(variableOutput().querySelector("em")).not.toBeNull();
  expect(variableOutput().innerHTML).toContain("updated html");
});
