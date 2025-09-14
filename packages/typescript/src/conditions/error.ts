import type { ObjectType } from "../datatypes/objects";

/**
 * @description
 * Forces a the typing to hard fail instead of using the soft fail mechanics of
 * the "never" type.
 * This is useful for quickly identifying type issues during development.
 *
 * @example
 * ```ts
 * type IsValid<T> = T extends string ? T : ErrorType<"Type is not a string">;
 * ```
 */
export type CompileError<
  Message extends string,
  // We use an object for context instead of a string template because some
  // types cannot be types cannot be represented inside of strings.
  // By using an Object, we can provide more detailed information such as
  // objects in the type error context.
  context extends ObjectType = ObjectType,
> = { error: Message; context: context };

/**
 * @description
 * An error thrown at runtime that can be caught during compile time.
 * If you do not remove the error type, you will receive a compile-time error.
 */
export type RuntimeError = Error;
