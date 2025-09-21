window.stateTracker = new WeakMap<WeakKey, unknown>();

/**
 * @description
 * A decorator that enables the transfer of state between server and client.
 */
export function hydrate<T>() {
  return (target: object, propertyKey: string | symbol): void => {
  };
}

declare global {
  interface Window {
    stateTracker: WeakMap<WeakKey, unknown>;
  }
}
