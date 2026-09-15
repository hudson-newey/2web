import { describe, expect, test } from "vitest";
import { normalizeDebugDocument } from "../src/debugDocument";

// The debug file is written by whatever compiler version built the inspected
// site, and older builds could carry partial entries over from even older
// builds. The normalizer must degrade any partial document to empty panels
// instead of crashing mid-render (which used to leave the panel stuck on the
// loading status).
describe("normalizeDebugDocument", () => {
  test("should reject non-object documents", () => {
    expect(normalizeDebugDocument(null)).toBeNull();
    expect(normalizeDebugDocument(undefined)).toBeNull();
    expect(normalizeDebugDocument("nope")).toBeNull();
    expect(normalizeDebugDocument(42)).toBeNull();
  });

  test("should fill missing top level arrays with empty lists", () => {
    const normalized = normalizeDebugDocument({ version: 2 })!;

    expect(normalized.version).toBe(2);
    expect(normalized.pages).toEqual([]);
    expect(normalized.assets).toEqual([]);
    expect(normalized.messages).toEqual([]);
    expect(normalized.reactivity).toEqual([]);
  });

  test("should fill missing per page arrays with empty lists", () => {
    const normalized = normalizeDebugDocument({
      version: 2,
      reactivity: [{ page: "index.html" }],
    })!;

    expect(normalized.reactivity.length).toBe(1);
    expect(normalized.reactivity[0].page).toBe("index.html");
    expect(normalized.reactivity[0].variables).toEqual([]);
    expect(normalized.reactivity[0].properties).toEqual([]);
    expect(normalized.reactivity[0].events).toEqual([]);
  });

  test("should skip null reactivity entries and null arrays", () => {
    const normalized = normalizeDebugDocument({
      version: 2,
      reactivity: [null, { page: "index.html", variables: null }],
    })!;

    expect(normalized.reactivity.length).toBe(1);
    expect(normalized.reactivity[0].variables).toEqual([]);
  });

  test("should fill missing entry fields with safe defaults", () => {
    const normalized = normalizeDebugDocument({
      version: 2,
      reactivity: [
        {
          page: "index.html",
          variables: [{ name: "$count" }],
          properties: [{ prop: "textContent" }],
          events: [{ event: "click" }],
        },
      ],
    })!;

    const [variable] = normalized.reactivity[0].variables;
    expect(variable).toEqual({
      name: "$count",
      initialValue: "",
      reactivity: "unknown",
      dependsOn: [],
    });

    const [property] = normalized.reactivity[0].properties;
    expect(property).toEqual({
      prop: "textContent",
      reducer: "",
      dependencies: [],
    });

    const [event] = normalized.reactivity[0].events;
    expect(event).toEqual({
      event: "click",
      reducer: "",
      sink: "",
      dependencies: [],
    });
  });

  test("should keep a well formed document intact", () => {
    const document = {
      version: 2,
      pages: ["index.html"],
      assets: [{ path: "index.js", type: "script", size: 12 }],
      messages: ["warning"],
      reactivity: [
        {
          page: "index.html",
          variables: [
            {
              name: "$count",
              initialValue: "0",
              reactivity: "reactive",
              dependsOn: ["$base"],
            },
          ],
          properties: [
            { prop: "textContent", reducer: "$count", dependencies: ["$count"] },
          ],
          events: [
            {
              event: "click",
              reducer: "$count = $count + 1",
              sink: "$count",
              dependencies: ["$count"],
            },
          ],
        },
      ],
    };

    expect(normalizeDebugDocument(document)).toEqual(document);
  });
});
