/**
 * @description
 * Excludes a type from another type.
 *
 * @example
 * ```ts
 * type NumbersWithoutInfinity = string & Not<Infinity>;
 */
export type Not<T, U> = T extends U ? never : T;
