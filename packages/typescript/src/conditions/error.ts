/**
 * @description
 * Forces a the typing to hard fail instead of using the soft fail mechanics of
 * the "never" type.
 * This is useful for quickly identifying type issues during development.
 *
 * @example
 * ```ts
 * type IsValid<T> = T extends string ? T : ErrorType<"Type is not a string">;
 */
export type ErrorType<Message extends string> = { error: Message };
