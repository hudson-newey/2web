import { beforeEach, expect, test, vi } from "vitest";
import { getByText } from "@testing-library/dom";
import userEvent from "@testing-library/user-event";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { BrowserFrame, Document } from "happy-dom";

// Verb route files ("__get.server.ts", "__post.server.ts", ...) are the only
// http entry points for server endpoints. Each one default exports an express
// request handler that is mounted at the url path that mirrors the directory
// the file sits in, for the http method the file name declares.

let document: Document;
let frame: BrowserFrame;
let fetchMock: ReturnType<typeof vi.fn>;

const getButton = () => getByText(document.body as any, "GET /api");
const postButton = () => getByText(document.body as any, "POST /api");

beforeEach(async () => {
  frame = await navigateToPage("verb-routes.html");
  document = frame.document;

  // The compiled functions call fetch at click time. happy-dom's fetch can't
  // reach the (not running) route server, so it is mocked out and the calls
  // are asserted instead of the responses.
  fetchMock = vi.fn(async () => ({ text: async () => "{}" }));
  (frame.window as any).fetch = fetchMock;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should call the get verb route at the directory mirrored path", async () => {
  const user = userEvent.setup();

  await user.click(getButton());
  await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));

  const [url, options] = fetchMock.mock.calls[0];

  expect(url).toBe("/api");

  // A default fetch() call has no init argument; its method is GET.
  expect(options?.method ?? "GET").toBe("GET");
});

test("should call the post verb route with a json body", async () => {
  const user = userEvent.setup();

  await user.click(postButton());
  await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1));

  const [url, options] = fetchMock.mock.calls[0];

  expect(url).toBe("/api");
  expect(options.method).toBe("POST");
  expect(options.headers["content-type"]).toBe("application/json");
  expect(JSON.parse(options.body)).toEqual({ from: "verb-routes.html" });
});
