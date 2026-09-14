import { beforeEach, expect, test, vi } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { BrowserFrame, Document } from "happy-dom";

// Functions declared in a <script compiled> block are passed through to the
// page's compiled JavaScript as normal functions, with two reactive
// extensions: assigning to a reactive variable inside the function body
// updates the page, and imported server functions can be called like normal
// async functions.

let document: Document;
let frame: BrowserFrame;
let fetchMock: ReturnType<typeof vi.fn>;

const increaseButton = () => getByText(document.body as any, "Increase");
const doubleButton = () => getByText(document.body as any, "Double");
const greetButton = () => getByText(document.body as any, "Greet");
const setUserButton = () => getByText(document.body as any, "Set user");
const countParagraph = () => getByText(document.body as any, "Count").closest("p")!;
const nameParagraph = () => getByText(document.body as any, "Name").closest("p")!;
const greetingParagraph = () => getByText(document.body as any, "Greeting").closest("p")!;

beforeEach(async () => {
  frame = await navigateToPage("functions.html");
  document = frame.document;

  // The rpc passthrough calls fetch at click time. happy-dom's fetch can't
  // reach the (not running) route server, so it is mocked out.
  fetchMock = vi.fn(async () => ({ text: async () => "Hello, world!" }));
  (frame.window as any).fetch = fetchMock;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should render the initial values", () => {
  expect(countParagraph().textContent).toContain("Count 1");
  expect(nameParagraph().textContent).toContain("Name world");
  expect(greetingParagraph().textContent).toContain("Greeting");
});

// A function that assigns to a reactive variable updates the page through the
// variable's update cascade.
test("should update reactive variables assigned in a function", async () => {
  const user = userEvent.setup();

  await user.click(increaseButton());
  expect(countParagraph().textContent).toContain("Count 2");

  await user.click(increaseButton());
  expect(countParagraph().textContent).toContain("Count 3");
});

test("should support compound assignments in a function", async () => {
  const user = userEvent.setup();

  // $count = 1 -> += $count -> 2
  await user.click(doubleButton());
  expect(countParagraph().textContent).toContain("Count 2");
});

test("should call imported server functions like normal async functions", async () => {
  const user = userEvent.setup();

  // A function assigned $name earlier, so the rpc call reads its CURRENT
  // runtime value.
  await user.click(setUserButton());
  expect(nameParagraph().textContent).toContain("Name 2web");

  await user.click(greetButton());
  await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));

  const [url, options] = fetchMock.mock.calls[0];

  expect(url).toBe("/_2web/rpc/api/greetings.server.js/greet");
  expect(JSON.parse(options.body)).toEqual(["2web"]);

  // The function's assignment to $greeting updates the page with the rpc
  // response.
  await vi.waitFor(() =>
    expect(greetingParagraph().textContent).toContain("Greeting Hello, world!"),
  );
});
