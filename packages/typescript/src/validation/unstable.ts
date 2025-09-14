import type { FunctionType, Synchronous } from "../datatypes/functions";
import type { Hide } from "../visibility/hide";

/**
 * @description
 * A type that cannot be narrowed down using type guards.
 * This is useful for unstable queries such as DOM queries or external data.
 *
 * ## Problem
 *
 * @example
 * ```ts
 * const element = document.querySelector(".my-element");
 * if (!element) {
 *   throw new Error("Element not found");
 * }
 *
 * setTimeout(() => {
 *   // This has the possibility to cause a runtime error because some other
 *   // part of the program might have removed this element from the DOM.
 *   // If you instead used the "Unstable" type, you would be forced to check
 *   // that the element still exists before using it.
 *   element.classList.add("active");
 * }, 1000);
 * ```
 *
 * ## Solution
 *
 * @example
 * ```ts
 * const element = unstable(document.querySelector(".my-element"));
 *
 * setTimeout(() => {
 *   // We are forced to use the "unwrap" function to type narrow the unstable
 *   // type to a stable type.
 *   unwrap(element, (stableElement) => {
 *     if (!stableElement) {
 *       throw new Error("Element not found");
 *     }
 *
 *     stableElement.classList.add("active");
 *   });
 * });
 * ```
 */
export type Unstable<T> = Hide<T>;

type Stable<T> = T;

/**
 * @description
 * Marks a value as "Unstable", preventing type narrowing outside of an `open`
 * function.
 */
export function unstable<T>(value: T): Unstable<T> {
  return value as Unstable<T>;
}

/**
 * @description
 * Unwraps a "Unstable" type into its stable counterpart without type narrowing
 * the unstable type.
 *
 * @example
 * ```ts
 * const element = Unstable<HTMLElement>(document.querySelector(".my-element"));
 *
 * setTimeout(() => {
 *   // We are forced to use the "unwrap" function to type narrow the unstable
 *   // type to a stable type.
 *   unwrap(element, (stableElement) => {
 *     if (!stableElement) {
 *       throw new Error("Element not found");
 *     }
 *
 *     stableElement.classList.add("active");
 *   });
 * });
 * ```
 */
export function unwrap<T, ReturnType>(
  ref: Unstable<T>,

  // We require the action to be synchronous because as soon as we defer the
  // task to the microtask queue, the value might have changed again.
  action: Synchronous<FunctionType<[ref: Stable<T>], ReturnType>>,
): ReturnType {
  return action(ref as Stable<T>);
}
