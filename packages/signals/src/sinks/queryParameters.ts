import type { Signal } from "../signal";

export function queryParameter<const Parameter extends string, T>(
  parameter: Parameter,
  signal: Signal<T>
) {
  signal.subscribe((value) => {
    if (typeof window === "undefined") {
      console.warn("queryParameter sink can only be used in a browser environment");
      return;
    }

    const url = new URL(window.location.href);
    if (value === undefined || value === null) {
      window.history.replaceState({}, "", url);
    } else {
      const stringifiedValue = value == null ? "" : String(value);
      url.searchParams.set(parameter, stringifiedValue);
    }

    window.history.replaceState({}, "", url);
  });
}
