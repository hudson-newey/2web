// @vitest-environment happy-dom
import { beforeEach, expect, test, vi } from "vitest";
import { Window } from "happy-dom";

// Integration test for the devtools panel: the panel's script runs against a
// mocked chrome.devtools API (like the real panel inside chrome) and the
// rendered panels are asserted. This reproduces the "0 variables, 0
// properties, 0 events" report: a debug document that loads must render its
// reactive graph, and an empty one must report zeros instead of crashing.

const panelDom = `
  <header>
    <button class="tab active" data-tab="assets">Assets</button>
    <button class="tab" data-tab="reactivity">Reactivity</button>
    <button class="tab" data-tab="graph">Graph</button>
    <button id="refresh">Refresh</button>
  </header>
  <main>
    <section id="panel-assets" class="panel active"></section>
    <section id="panel-reactivity" class="panel"></section>
    <section id="panel-graph" class="panel"></section>
  </main>
  <footer id="status"></footer>
`;

// A representative __2web.debug.json document (the shape the compiler writes
// for a development build).
const debugDocument = {
  version: 2,
  pages: ["mouse-position-demo.html"],
  assets: [{ path: "mouse-position-demo.html", type: "page", size: 1024 }],
  messages: [],
  reactivity: [
    {
      page: "dev/mouse-position-demo.html",
      variables: [
        { name: "$mouseX", initialValue: "0", reactivity: "assignment", dependsOn: null },
        { name: "$mouseY", initialValue: "0", reactivity: "assignment", dependsOn: null },
      ],
      properties: [
        { prop: "textContent", reducer: " $mouseX ", dependencies: ["$mouseX"] },
        { prop: "textContent", reducer: " $mouseY ", dependencies: ["$mouseY"] },
      ],
      events: [
        { event: "mousemove", reducer: "$mouseX = event.x", sink: "$mouseX", dependencies: [] },
      ],
    },
  ],
};

type EvalCallback = (result: unknown, exceptionInfo?: unknown) => void;

let evalResponse: unknown;

function mockDevtoolsApi(evalCallback: EvalCallback): void {
  (globalThis as any).chrome = {
    devtools: {
      inspectedWindow: {
        eval: (_expression: string, callback: EvalCallback) => evalCallback(callback),
      },
      network: {
        onNavigated: { addListener: (_listener: () => void) => {} },
      },
    },
  };
}

async function loadPanel(): Promise<Window["document"]> {
  document.body.innerHTML = panelDom;

  await import("../src/script.ts");

  // The panel refreshes asynchronously on load.
  await vi.waitFor(() => {
    expect(document.getElementById("status")!.textContent).not.toBe("");
  });

  return document;
}

beforeEach(() => {
  vi.resetModules();
  evalResponse = debugDocument;
  mockDevtoolsApi((callback) => callback(evalResponse));
});

test("should render the reactive graph of a debug document", async () => {
  const document = await loadPanel();

  expect(document.getElementById("status")!.textContent).toContain("2 variables");
  expect(document.getElementById("status")!.textContent).toContain("2 properties");
  expect(document.getElementById("status")!.textContent).toContain("1 events");

  // The reactivity panel renders one table row per variable.
  const variableRows = document.querySelectorAll("#panel-reactivity table")[0].querySelectorAll("tbody tr, tr");
  expect(variableRows.length).toBeGreaterThanOrEqual(2);

  // The graph panel renders one node per graph entry.
  const graphNodes = document.querySelectorAll("#panel-graph .graph-node");
  expect(graphNodes.length).toBe(5);
});

test("should report zero counts for an empty reactive graph", async () => {
  evalResponse = { ...debugDocument, reactivity: [] };

  const document = await loadPanel();

  expect(document.getElementById("status")!.textContent).toContain("0 variables");
  expect(document.getElementById("status")!.textContent).toContain("0 properties");
  expect(document.getElementById("status")!.textContent).toContain("0 events");
});

test("should render partial (older format) debug documents without crashing", async () => {
  // Entries written by older compilers can miss fields entirely.
  evalResponse = {
    version: 1,
    pages: null,
    assets: null,
    messages: null,
    reactivity: [
      { page: "index.html", variables: null, properties: null, events: null },
      null,
    ],
  };

  const document = await loadPanel();

  expect(document.getElementById("status")!.textContent).toContain("0 variables");
  expect(document.getElementById("panel-reactivity")!.textContent).not.toContain("$");
});

test("should report a failed inspection instead of hanging on the loading status", async () => {
  mockDevtoolsApi((callback) =>
    callback(undefined, {
      isError: true,
      code: "E_NOTFOUND",
      description: "inspected page is restricted",
      details: [],
      isException: false,
      value: "",
    }),
  );

  const document = await loadPanel();

  expect(document.getElementById("status")!.textContent).toContain("not found");
});

test("should refresh the panels when the refresh button is clicked", async () => {
  const document = await loadPanel();

  evalResponse = { ...debugDocument, reactivity: [] };

  document.getElementById("refresh")!.click();

  await vi.waitFor(() => {
    expect(document.getElementById("status")!.textContent).toContain("0 variables");
  });
});

test("should explain that the devtools api is unavailable outside chrome", async () => {
  vi.resetModules();
  (globalThis as any).chrome = undefined;

  const document = await loadPanel();

  expect(document.getElementById("status")!.textContent).toContain("chrome.devtools API unavailable");
});
