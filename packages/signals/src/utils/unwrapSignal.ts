import { onFirstUpdate } from "../lifecycle/onFirstUpdate";
import type { Signal } from "../signal";
import { isSignal } from "./isSignal";

export type MaybeSignal<T> = T | Signal<T> | Promise<T> | Promise<Signal<T>>;

export async function unwrapSignal<const T>(value: MaybeSignal<T>): Promise<T> {
  const resValue = await value;
  if (isSignal(resValue)) {
    await onFirstUpdate(resValue);
    return resValue.value as T;
  }

  return resValue as T;
}
