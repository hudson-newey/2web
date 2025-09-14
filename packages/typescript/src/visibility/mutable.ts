/**
 * @description
 * Removes the `readonly` modifier from all properties in an object type `T`.
 */
export type Mutable<T> = {
  -readonly [Property in keyof T]: T[Property];
};

export function mutable<const T>(value: T): Mutable<T> {
  return value as Mutable<T>;
}
