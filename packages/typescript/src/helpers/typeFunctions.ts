/**
 * A JavaScript function that can is invoked at compile time.
 * This can be used to perform complex type manipulations that would otherwise
 * require long and unwieldy type definitions.
 *
 * @example
 * ```ts
 * type NoNull<T> =
 *      typeof typeFunction((value: T) => {
 *        if (value === null) {
 *          throw new Error("Null values are not allowed");
 *        }
 *
 *        return value;
 *      });
 * ```
 */
export function typeFunction() {}
