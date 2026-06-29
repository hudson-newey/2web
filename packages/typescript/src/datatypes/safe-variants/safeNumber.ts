/**
 * Enforces that your number datatype cannot be accidently type narrowed using
 * `0` as a falsy value.
 *
 * @example
 * ```ts
 * const value: SafeNumber | null = 0;
 *
 * // The following will fail during compilation because `value` was incorrectly
 * // narrowed using a falsy condition.
 * if (!value) throw new Error("value must be provided");
 *
 * // The following would be the expected behavior.
 * if (value === null) throw new Error("value must be provided");
 * ````
 */
export type SafeNumber<T extends number> = Omit<T, 0> | 0;
