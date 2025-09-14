/**
 * @description
 * A helper type that removes keys `U` from object type `T`.
 * This is similar to `Omit` and `Exclude`, but it can automatically infer what
 * omission to perform based on the types.
 */
export type Without<T, U> = U extends string | number | symbol
  ? Omit<T, U>
  : Exclude<T, U>;
