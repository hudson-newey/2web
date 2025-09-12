import { Signal } from "../signal";

export function isSignal(value: any): value is Signal<unknown> {
  return value instanceof Signal;
}
