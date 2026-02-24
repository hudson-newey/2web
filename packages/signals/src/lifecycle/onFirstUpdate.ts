import type { Signal } from "../signal";

// TODO: Make this function act as an "asserts" function that removes the
// uninitialized value from a signal.
/**
 * @description
 * Hooks into the onCreate lifecycle of a signal and returns a promise that will
 * resolve when the onCreate lifecycle is triggered.
 *
 * @returns
 * The signal that triggered the onCreate lifecycle.
 * The signal is returned so that you can use this lifecycle hook in the
 * creation of a signal or use helpers like {@linkcode Promise.race} and
 * {@linkcode Promise.any}.
 */
export function onFirstUpdate<T extends Signal<U>, U = any>(
  signal: Readonly<T>,
): Promise<T> {
  // If the signal has already completed the onFirstUpdate lifecycle,
  // resolve immediately.
  //
  // Note that doneFirstUpdate is an internal property and is not public so we
  // therefore need to cast to "any" to access it.
  if ((signal as any).doneFirstUpdate) {
    return Promise.resolve(signal);
  }

  return new Promise((res) => {
    signal.onFirstUpdate(() => {
      res(signal);
    });
  });
}

export type OnFirstUpdateCallback<T> = (value: T | null) => void;
