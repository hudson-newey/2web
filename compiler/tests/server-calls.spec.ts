import { beforeEach, expect, test, vi } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { BrowserFrame, Document } from "happy-dom";

// Importing a server script (.server.ts) from a <script compiled> block
// generates client side rpc passthroughs for the imported functions. A
// reactive event listener that calls the imported function directly compiles
// into a listener that posts the call arguments (as a json array) to the
// generated rpc endpoint.

let document: Document;
let frame: BrowserFrame;
let fetchMock: ReturnType<typeof vi.fn>;

const greetButton = () => getByText(document.body as any, "Greet");
const directButton = () => getByText(document.body as any, "Greet direct");

beforeEach(async () => {
  frame = await navigateToPage("server-calls.html");
  document = frame.document;

  // The rpc passthrough calls fetch at click time. happy-dom's fetch can't
  // reach the (not running) route server, so it is mocked out and the calls
  // are asserted instead of the responses.
  fetchMock = vi.fn(async () => ({ text: async () => "Hello, 2web!" }));
  (frame.window as any).fetch = fetchMock;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should call the rpc endpoint of the imported server function with the reactive state", async () => {
  const user = userEvent.setup();

  await user.click(greetButton());
  await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));

  const [url, options] = fetchMock.mock.calls[0];

  expect(url).toBe("/_2web/rpc/api/greetings.server.js/greet");
  expect(options.method).toBe("POST");
  expect(options.headers["content-type"]).toBe("application/json");

  // Static reactive variables are inlined into the call arguments.
  expect(JSON.parse(options.body)).toEqual(["2web"]);
});

test("should pass literal call arguments through to the rpc endpoint", async () => {
  const user = userEvent.setup();

  await user.click(directButton());
  await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));

  expect(fetchMock.mock.calls[0][0]).toBe(
    "/_2web/rpc/api/greetings.server.js/greet",
  );
  expect(JSON.parse(fetchMock.mock.calls[0][1].body)).toEqual(["direct"]);
});
