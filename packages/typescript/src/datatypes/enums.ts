export type Enum<T extends string> = { [K in T]: symbol };

/**
 * @description
 * JavaScript compatible enum creation function.
 */
export function createEnum<const T extends string>(keys: T[]): Enum<T> {
  const enumObj = {} as Enum<T>;

  for (const v of keys) {
    enumObj[v] = Symbol(v);
  }

  return Object.freeze(enumObj);
}
