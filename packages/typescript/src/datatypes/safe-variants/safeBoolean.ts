/**
 * Enforces that your boolean datatype cannot be accidently type narrowed using
 * a falsy check.
 *
 * @example
 * ```ts
 * const value: SafeBoolean | null = false;
 *
 * // The following will fail during compilation because `value` was incorrectly
 * // narrowed using a falsy condition.
 * if (!value) throw new Error("value must be provided");
 *
 * // The following would be the expected behavior.
 * if (value === null) throw new Error("value must be provided");
 * ````
 */
export type SafeBoolean<T extends boolean> = Exclude<T, false> | false;
