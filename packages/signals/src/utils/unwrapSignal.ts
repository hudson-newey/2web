import type { Signal } from "../signal";
import { isSignal } from "./isSignal";

export type MaybeSignal<T> = T | Signal<T>;

export function unwrapSignal<const T>(value: MaybeSignal<T>): T {
  if (isSignal(value)) {
    return value.value as T;
  }

  return value as T;
}
