import { Signal } from "../signal";

export function isSignal<T>(value: unknown | Signal<T>): value is Signal<T> {
  return value instanceof Signal;
}
