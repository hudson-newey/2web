import type { CompileError } from "./error";

/**
 * @description
 * A conditional type tha can be used to void a type intersection based on a
 * boolean condition.
 *
 * @example
 * ```ts
 * type WithToString<T> = T & If<T extends { toString(): string }>;
 * ```
 */
export type If<Type, ExpectedType> = Type extends ExpectedType
  ? unknown
  : CompileError<`Type condition failed`>;
