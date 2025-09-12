import { ReadonlySignal } from "../readonlySignal";

export function isReadonlySignal(value: any): value is ReadonlySignal<unknown> {
  return value instanceof ReadonlySignal;
}
