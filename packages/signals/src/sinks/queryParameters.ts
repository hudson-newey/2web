import type { Signal } from "../signal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export function queryParameter<const Parameter extends string, T>(
  parameter: MaybeSignal<Parameter>,
  signal: Signal<T>
) {
  signal.subscribe((value) => {
    if (typeof window === "undefined") {
      console.warn(
        "queryParameter sink can only be used in a browser environment"
      );
      return;
    }

    const url = new URL(window.location.href);
    if (value === undefined || value === null) {
      window.history.replaceState({}, "", url);
    } else {
      const stringifiedValue = value == null ? "" : String(value);
      url.searchParams.set(unwrapSignal(parameter), stringifiedValue);
    }

    window.history.replaceState({}, "", url);
  });
}
