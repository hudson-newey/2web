import type { Signal } from "../signal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export async function queryParameter<const Parameter extends string, T>(
  parameter: MaybeSignal<Parameter>,
  signal: Signal<T>,
) {
  const fn = async (value: T | null) => {
    if (typeof window === "undefined") {
      console.warn(
        "queryParameter sink can only be used in a browser environment",
      );
      return;
    }

    const url = new URL(window.location.href);
    if (value === undefined || value === null) {
      window.history.replaceState({}, "", url);
    } else {
      const stringifiedValue = value == null ? "" : String(value);
      const param = await unwrapSignal(parameter);

      url.searchParams.set(param, stringifiedValue);
    }

    window.history.replaceState({}, "", url);
  };

  signal.subscribe(fn);

  if (isSignal(parameter)) {
    parameter.subscribe(() => fn(signal.value));
  }
}
