// The debug file (__2web.debug.json) is written by whatever compiler version
// built the inspected site, and (before the compiler pruned old entries) it
// could contain partial entries carried over from older builds. Every array
// and field of the document is therefore coerced into the shape the panel
// renders, so that a partial or stale debug file degrades to empty panels
// instead of crashing mid-render.

export interface VariableInfo {
  name: string;
  initialValue: string;
  reactivity: string;
  dependsOn?: string[];
}

export interface PropertyInfo {
  prop: string;
  reducer: string;
  dependencies: string[];
}

export interface EventInfo {
  event: string;
  reducer: string;
  sink: string;
  dependencies: string[];
}

export interface PageReactivity {
  page: string;
  variables: VariableInfo[];
  properties: PropertyInfo[];
  events: EventInfo[];
}

export interface AssetInfo {
  path: string;
  type: string;
  size: number;
}

export interface DebugDocument {
  version: number;
  pages: string[];
  assets: AssetInfo[];
  messages: string[];
  reactivity: PageReactivity[];
}

export function normalizeDebugDocument(raw: unknown): DebugDocument | null {
  if (!raw || typeof raw !== "object") {
    return null;
  }

  const document_ = raw as Record<string, unknown>;

  const list = (value: unknown): unknown[] => (Array.isArray(value) ? value : []);
  const text = (value: unknown, fallback: string): string =>
    typeof value === "string" ? value : fallback;

  const reactivity = list(document_.reactivity)
    .filter((entry) => entry && typeof entry === "object")
    .map((entry) => {
      const page = entry as Record<string, unknown>;

      const variables = list(page.variables).map((variable) => {
        const info = (variable ?? {}) as Record<string, unknown>;

        return {
          name: text(info.name, "?"),
          initialValue: text(info.initialValue, ""),
          reactivity: text(info.reactivity, "unknown"),
          dependsOn: list(info.dependsOn) as string[],
        };
      });

      const properties = list(page.properties).map((property) => {
        const info = (property ?? {}) as Record<string, unknown>;

        return {
          prop: text(info.prop, "?"),
          reducer: text(info.reducer, ""),
          dependencies: list(info.dependencies) as string[],
        };
      });

      const events = list(page.events).map((event) => {
        const info = (event ?? {}) as Record<string, unknown>;

        return {
          event: text(info.event, "?"),
          reducer: text(info.reducer, ""),
          sink: text(info.sink, ""),
          dependencies: list(info.dependencies) as string[],
        };
      });

      return {
        page: text(page.page, "unknown page"),
        variables,
        properties,
        events,
      };
    });

  return {
    version: typeof document_.version === "number" ? document_.version : 0,
    pages: list(document_.pages) as string[],
    assets: list(document_.assets) as AssetInfo[],
    messages: list(document_.messages) as string[],
    reactivity,
  };
}
