import type { Signal } from "../signal";

export function onDestroy<T extends Signal<U>, U = any>(
  signal: Readonly<T>,
): Promise<T> {
  return new Promise((res) => {
    signal.onDestroy(() => {
      res(signal);
    });
  });
}

export type OnDestroyCallback<T> = (value: T | null) => void;
