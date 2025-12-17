/**
 * @description
 * A value that will automatically be zeroed out after use.
 */
export function ephemeral<const T>(value: T): T {
  const valueRef = [value];

  const resource = {
    [Symbol.dispose]() {
      valueRef[0] = undefined as unknown as T;
    },
  };

  // Create a proxy object for the resource so that
  return new Proxy(resource, {
    get() {
      return valueRef[0];
    },
  }) as unknown as T;
}
