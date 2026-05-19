import type { CompileError } from "./error";

/**
 * @description
 * A conditional type that can be used to void a type intersection based on a
 * boolean condition.
 *
 * @example
 * ```ts
 * type WithToString<T> = T & If<T, { toString(): string }>;
 * ```
 */
export type If<Type, Expected> = Type extends Expected
  ? unknown
  : CompileError<`Type condition failed`>;
