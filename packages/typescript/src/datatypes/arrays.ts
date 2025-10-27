/**
 * @description
 * An array where all elements are unique.
 * This is similar to a Set, but only exists at compile time, meaning that there
 * is no runtime overhead.
 * Because the underlying datatype is a regular array, all array methods are
 * available.
 */
export type UniqueArray<T extends readonly unknown[]> =
  T extends readonly [infer Head, ...infer Tail]
    ? Head extends Tail[number]
      ? never
      : readonly [Head, ...UniqueArray<Tail>]
    : T;
