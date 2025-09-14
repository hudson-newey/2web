import { ReadonlySignal } from "../readonlySignal";

export function isReadonlySignal<T>(
  value: unknown | ReadonlySignal<T>,
): value is ReadonlySignal<T> {
  return value instanceof ReadonlySignal;
}
