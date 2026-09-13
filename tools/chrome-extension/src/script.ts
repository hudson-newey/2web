// 2web devtools panel.
//
// The panel fetches the __2web.debug.json file that the 2web compiler writes
// into the build output of development builds and renders three views:
//
//   - Assets: every file in the compiled output
//   - Reactivity: the reactive variables, properties, and events of the site
//   - Graph: the relationship graph between variables, properties, and events

interface VariableInfo {
  name: string;
  initialValue: string;
  reactivity: string;
  dependsOn?: string[];
}

interface PropertyInfo {
  prop: string;
  reducer: string;
  dependencies: string[];
}

interface EventInfo {
  event: string;
  reducer: string;
  sink: string;
  dependencies: string[];
}

interface PageReactivity {
  page: string;
  variables: VariableInfo[];
  properties: PropertyInfo[];
  events: EventInfo[];
}

interface AssetInfo {
  path: string;
  type: string;
  size: number;
}

interface DebugDocument {
  version: number;
  pages: string[];
  assets: AssetInfo[];
  messages: string[];
  reactivity: PageReactivity[];
}

const DEBUG_FILE = "__2web.debug.json";

const STATUS = document.getElementById("status")!;
const PANELS: Record<string, HTMLElement> = {
  assets: document.getElementById("panel-assets")!,
  reactivity: document.getElementById("panel-reactivity")!,
  graph: document.getElementById("panel-graph")!,
};

// ---------------------------------------------------------------------------
// Data loading
// ---------------------------------------------------------------------------

// The debug file lives next to the compiled page, so it is fetched from the
// inspected page's origin through inspectedWindow.eval. This doesn't require
// any manifest permissions and works regardless of whether the debug file was
// requested during the page load.
async function fetchDebugDocument(): Promise<DebugDocument | null> {
  const expression = `
    fetch("/${DEBUG_FILE}")
      .then((response) => response.ok ? response.text() : Promise.reject(new Error("not found")))
      .then((body) => JSON.parse(body))
      .catch(() => null)
  `;

  return new Promise((resolve) => {
    chrome.devtools.inspectedWindow.eval(
      expression,
      (result: DebugDocument | null) => resolve(result ?? null),
    );
  });
}

function setStatus(message: string, isError = false): void {
  STATUS.textContent = message;
  STATUS.classList.toggle("error", isError);
}

// ---------------------------------------------------------------------------
// DOM helpers
// ---------------------------------------------------------------------------

function element<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  className?: string,
  textContent?: string,
): HTMLElementTagNameMap[K] {
  const node = document.createElement(tag);
  if (className) {
    node.className = className;
  }
  if (textContent !== undefined) {
    node.textContent = textContent;
  }
  return node;
}

function heading(text: string): HTMLElement {
  return element("h2", "section-heading", text);
}

function table(headers: string[], rows: HTMLElement[][]): HTMLElement {
  const tableNode = element("table", "data-table");
  const headRow = element("tr");

  for (const header of headers) {
    headRow.appendChild(element("th", undefined, header));
  }
  tableNode.appendChild(headRow);

  for (const row of rows) {
    const rowNode = element("tr");
    for (const cell of row) {
      rowNode.appendChild(cell);
    }
    tableNode.appendChild(rowNode);
  }

  return tableNode;
}

function cell(text: string, className?: string): HTMLElement {
  return element("td", className, text);
}

function emptyState(message: string): HTMLElement {
  return element("p", "empty-state", message);
}

// ---------------------------------------------------------------------------
// Assets tab
// ---------------------------------------------------------------------------

function formatBytes(size: number): string {
  if (size < 1024) {
    return `${size} B`;
  }

  if (size < 1024 * 1024) {
    return `${(size / 1024).toFixed(1)} KiB`;
  }

  return `${(size / (1024 * 1024)).toFixed(1)} MiB`;
}

function renderAssets(panel: HTMLElement, document_: DebugDocument): void {
  panel.replaceChildren();

  panel.appendChild(
    heading(`Assets (${document_.assets.length}) · Pages (${document_.pages.length})`),
  );

  if (document_.assets.length === 0) {
    panel.appendChild(emptyState("No build output found."));
    return;
  }

  const rows = document_.assets.map((asset) => [
    cell(asset.path, "mono"),
    cell(asset.type, `asset-type asset-${asset.type}`),
    cell(formatBytes(asset.size), "numeric"),
  ]);

  panel.appendChild(table(["Path", "Type", "Size"], rows));
}

// ---------------------------------------------------------------------------
// Reactivity tab
// ---------------------------------------------------------------------------

function renderReactivity(panel: HTMLElement, document_: DebugDocument): void {
  panel.replaceChildren();

  const variables: VariableInfo[] = [];
  const properties: PropertyInfo[] = [];
  const events: EventInfo[] = [];

  for (const pageGraph of document_.reactivity) {
    variables.push(...pageGraph.variables);
    properties.push(...pageGraph.properties);
    events.push(...pageGraph.events);
  }

  panel.appendChild(heading(`Variables (${variables.length})`));
  if (variables.length === 0) {
    panel.appendChild(emptyState("No reactive variables were compiled."));
  } else {
    panel.appendChild(
      table(
        ["Variable", "Initial value", "Reactivity"],
        variables.map((variable) => [
          cell(variable.name, "mono variable-name"),
          cell(variable.initialValue, "mono"),
          cell(variable.reactivity, `reactivity-badge reactivity-${variable.reactivity}`),
        ]),
      ),
    );
  }

  panel.appendChild(heading(`Properties (${properties.length})`));
  if (properties.length === 0) {
    panel.appendChild(emptyState("No reactive properties were compiled."));
  } else {
    panel.appendChild(
      table(
        ["Property", "Reducer", "Dependencies"],
        properties.map((property) => [
          cell(`*${property.prop}`, "mono"),
          cell(property.reducer, "mono"),
          cell(property.dependencies.join(", "), "mono"),
        ]),
      ),
    );
  }

  panel.appendChild(heading(`Events (${events.length})`));
  if (events.length === 0) {
    panel.appendChild(emptyState("No reactive events were compiled."));
  } else {
    panel.appendChild(
      table(
        ["Event", "Reducer", "Sink", "Dependencies"],
        events.map((event) => [
          cell(`@${event.event}`, "mono"),
          cell(event.reducer, "mono"),
          cell(event.sink || "—", "mono"),
          cell(event.dependencies.join(", "), "mono"),
        ]),
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Graph tab
// ---------------------------------------------------------------------------

interface GraphNode {
  id: string;
  label: string;
  kind: "variable" | "property" | "event";
  page: string;
  x: number;
  y: number;
}

interface GraphEdge {
  from: string;
  to: string;
  kind: "sink" | "dependency";
}

const SVG_NS = "http://www.w3.org/2000/svg";
const COLUMN_X: Record<GraphNode["kind"], number> = {
  variable: 140,
  property: 420,
  event: 700,
};
const ROW_HEIGHT = 34;
const TOP_MARGIN = 30;
const NODE_WIDTH = 150;

function buildGraph(document_: DebugDocument): { nodes: GraphNode[]; edges: GraphEdge[] } {
  const nodes = new Map<string, GraphNode>();
  const edges: GraphEdge[] = [];

  const ensureNode = (id: string, label: string, kind: GraphNode["kind"], page: string): GraphNode => {
    let node = nodes.get(id);
    if (!node) {
      node = { id, label, kind, page, x: 0, y: 0 };
      nodes.set(id, node);
    }
    return node;
  };

  for (const pageGraph of document_.reactivity) {
    const source = pageGraph.page;

    for (const variable of pageGraph.variables) {
      ensureNode(variable.name, variable.name, "variable", source);
    }

    // Variable-to-variable edges: a computed variable is redrawn whenever the
    // variables it is computed from change.
    for (const variable of pageGraph.variables) {
      for (const dependency of variable.dependsOn ?? []) {
        if (nodes.has(dependency) && nodes.has(variable.name)) {
          edges.push({ from: dependency, to: variable.name, kind: "dependency" });
        }
      }
    }

    for (const property of pageGraph.properties) {
      const id = `prop:${source}:${property.prop}:${property.reducer}`;
      ensureNode(id, `*${property.prop} = ${property.reducer.trim()}`, "property", source);

      for (const dependency of property.dependencies) {
        if (nodes.has(dependency)) {
          edges.push({ from: id, to: dependency, kind: "dependency" });
        }
      }
    }

    for (const event of pageGraph.events) {
      const id = `event:${source}:${event.reducer}`;
      ensureNode(id, `@${event.event}: ${event.reducer.trim()}`, "event", source);

      if (event.sink && nodes.has(event.sink)) {
        edges.push({ from: id, to: event.sink, kind: "sink" });
      }

      for (const dependency of event.dependencies) {
        if (nodes.has(dependency) && dependency !== event.sink) {
          edges.push({ from: id, to: dependency, kind: "dependency" });
        }
      }
    }
  }

  return { nodes: [...nodes.values()], edges };
}

// Nodes are laid out in three vertical columns (variables, properties, events)
// ordered by the page they were compiled from so that related nodes sit close
// together.
function layoutGraph(nodes: GraphNode[]): number {
  const nodesByKind: Record<GraphNode["kind"], GraphNode[]> = {
    variable: [],
    property: [],
    event: [],
  };

  for (const node of nodes) {
    nodesByKind[node.kind].push(node);
  }

  let height = 0;
  for (const kind of Object.keys(nodesByKind) as GraphNode["kind"][]) {
    const column = nodesByKind[kind].sort((a, b) => a.page.localeCompare(b.page) || a.label.localeCompare(b.label));

    column.forEach((node, index) => {
      node.x = COLUMN_X[kind];
      node.y = TOP_MARGIN + index * ROW_HEIGHT;
    });

    height = Math.max(height, TOP_MARGIN + column.length * ROW_HEIGHT);
  }

  return height;
}

function renderGraph(panel: HTMLElement, document_: DebugDocument): void {
  panel.replaceChildren();
  panel.appendChild(heading("Reactive relationships"));

  const { nodes, edges } = buildGraph(document_);

  if (nodes.length === 0) {
    panel.appendChild(emptyState("No reactive graph to display."));
    return;
  }

  const graphHeight = layoutGraph(nodes);

  const svg = document.createElementNS(SVG_NS, "svg");
  svg.setAttribute("width", "100%");
  svg.setAttribute("height", String(graphHeight + TOP_MARGIN));
  svg.classList.add("graph");

  const edgeLayer = document.createElementNS(SVG_NS, "g");
  const nodeLayer = document.createElementNS(SVG_NS, "g");

  for (const edge of edges) {
    const from = nodes.find((node) => node.id === edge.from)!;
    const to = nodes.find((node) => node.id === edge.to)!;

    const startX = from.x + NODE_WIDTH / 2;
    const endX = to.x - NODE_WIDTH / 2;
    const midX = (startX + endX) / 2;

    const path = document.createElementNS(SVG_NS, "path");
    path.setAttribute(
      "d",
      `M ${startX} ${from.y} C ${midX} ${from.y}, ${midX} ${to.y}, ${endX} ${to.y}`,
    );
    path.classList.add("edge", `edge-${edge.kind}`);
    path.dataset.from = edge.from;
    path.dataset.to = edge.to;

    edgeLayer.appendChild(path);
  }

  for (const node of nodes) {
    const group = document.createElementNS(SVG_NS, "g");
    group.classList.add("graph-node", `node-${node.kind}`);
    group.dataset.id = node.id;
    group.dataset.page = node.page;

    const rect = document.createElementNS(SVG_NS, "rect");
    rect.setAttribute("x", String(node.x - NODE_WIDTH / 2));
    rect.setAttribute("y", String(node.y - 12));
    rect.setAttribute("width", String(NODE_WIDTH));
    rect.setAttribute("height", "24");
    rect.setAttribute("rx", "6");

    const label = document.createElementNS(SVG_NS, "text");
    label.setAttribute("x", String(node.x));
    label.setAttribute("y", String(node.y + 4));
    label.setAttribute("text-anchor", "middle");
    label.textContent = node.label;

    group.appendChild(rect);
    group.appendChild(label);
    group.appendChild(document.createElementNS(SVG_NS, "title")).textContent = node.page;

    group.addEventListener("click", () => highlightConnections(group, svg));
    nodeLayer.appendChild(group);
  }

  svg.appendChild(edgeLayer);
  svg.appendChild(nodeLayer);
  panel.appendChild(svg);
}

function highlightConnections(selected: SVGGElement, svg: SVGSVGElement): void {
  const alreadySelected = selected.classList.contains("selected");
  svg.classList.remove("has-selection");

  for (const node of svg.querySelectorAll(".graph-node")) {
    node.classList.remove("selected", "connected");
  }
  for (const edge of svg.querySelectorAll(".edge")) {
    edge.classList.remove("connected");
  }

  if (alreadySelected) {
    return;
  }

  const selectedId = selected.dataset.id!;
  svg.classList.add("has-selection");
  selected.classList.add("selected");

  for (const edge of svg.querySelectorAll(".edge")) {
    const edgeElement = edge as SVGPathElement;
    const touches =
      edgeElement.dataset.from === selectedId || edgeElement.dataset.to === selectedId;
    if (!touches) {
      continue;
    }

    edge.classList.add("connected");

    const connectedId =
      edgeElement.dataset.from === selectedId ? edgeElement.dataset.to! : edgeElement.dataset.from!;
    const connectedNode = svg.querySelector(`.graph-node[data-id="${CSS.escape(connectedId)}"]`);
    connectedNode?.classList.add("connected");
  }
}

// ---------------------------------------------------------------------------
// Wiring
// ---------------------------------------------------------------------------

function switchTab(tabName: string): void {
  for (const tab of document.querySelectorAll(".tab")) {
    tab.classList.toggle("active", (tab as HTMLElement).dataset.tab === tabName);
  }

  for (const [name, panel] of Object.entries(PANELS)) {
    panel.classList.toggle("active", name === tabName);
  }
}

async function refresh(): Promise<void> {
  setStatus("Loading __2web.debug.json …");

  const debugDocument = await fetchDebugDocument();
  if (!debugDocument) {
    for (const panel of Object.values(PANELS)) {
      panel.replaceChildren();
      panel.appendChild(
        emptyState(
          `Could not load /${DEBUG_FILE}. Make sure the inspected page is a 2web development build (debug info is omitted from production builds).`,
        ),
      );
    }
    setStatus(`/${DEBUG_FILE} not found`, true);
    return;
  }

  renderAssets(PANELS.assets, debugDocument);
  renderReactivity(PANELS.reactivity, debugDocument);
  renderGraph(PANELS.graph, debugDocument);

  const variableCount = debugDocument.reactivity.reduce((sum, page) => sum + page.variables.length, 0);
  const eventCount = debugDocument.reactivity.reduce((sum, page) => sum + page.events.length, 0);
  const propertyCount = debugDocument.reactivity.reduce((sum, page) => sum + page.properties.length, 0);

  setStatus(
    `__2web.debug.json v${debugDocument.version} — ${debugDocument.assets.length} assets, ` +
      `${variableCount} variables, ${propertyCount} properties, ${eventCount} events`,
  );
}

for (const tab of document.querySelectorAll(".tab")) {
  tab.addEventListener("click", () => switchTab((tab as HTMLElement).dataset.tab!));
}

document.getElementById("refresh")!.addEventListener("click", refresh);
chrome.devtools.network.onNavigated.addListener(refresh);

refresh();
