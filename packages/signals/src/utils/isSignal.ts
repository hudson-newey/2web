import { Signal } from "../signal";

export function isSignal<const T>(
  value: T | Promise<T> | Signal<T> | Promise<Signal<T>>
): value is Signal<T> {
  return value instanceof Signal;
}
