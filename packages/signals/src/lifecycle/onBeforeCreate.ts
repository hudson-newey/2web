import type { Signal } from "../signal";

export type OnCreateCallback = () => void;

export function onCreate<T extends Signal<U>, U = any>(
  signal: Readonly<T>,
): Promise<T> {
  // If the signal has already completed the onCreate lifecycle, resolve
  // immediately.
  //
  // Note that doneCreate is an internal property and is not public so we
  // therefore need to cast to "any" to access it.
  if ((signal as any).doneCreate) {
    return Promise.resolve(signal);
  }

  return new Promise((res) => {
    signal.onCreate(() => {
      res(signal);
    });
  });
}
