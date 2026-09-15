/**
 * Compile time virtual functions.
 *
 * These functions don't exist at runtime: the compiler evaluates them while
 * it compiles a page. Import them in a `<script compiled>` block like a
 * normal esm import:
 *
 * ```html
 * <script compiled>
 *   import { $uid, $props, $env, $readFile } from "@compiler/types/interop.types.ts";
 * </script>
 * ```
 */

/**
 * Expands to a unique incrementing id (e.g. "uid-1"), which is useful for
 * aria attributes:
 *
 * ```html
 * <input aria-describedby="$uid()" />
 * ```
 *
 * The ids increment in document order and are deterministic per page. The
 * expansion is a bare value, so when you use it inside a script block, wrap
 * it in a string literal: `const id = "$uid()"`.
 */
export declare function $uid(): string;

/**
 * Declares the parameters a component takes. Components receive parameters
 * through bracket attributes on an instance of the component:
 *
 * ```html
 * <Counter [title]="'My Counter'" [count]="1" />
 * ```
 *
 * `$props().title` expands to the value the instance passed for `title`.
 * Reactive variables stay reactive: the component updates when the passed
 * reactive variable changes.
 */
export declare function $props(): Record<string, any>;

/**
 * Expands to the value of an environment variable (as a string literal) at
 * compile time.
 *
 * The value is a compile time constant, not reactive. Assign it to a reactive
 * variable to render it.
 */
export declare function $env(name: string): string;

/**
 * Expands to the content of a file (as a string literal) at compile time.
 * Paths are resolved relative to the file that reads it.
 *
 * The value is a compile time constant, not reactive. Assign it to a reactive
 * variable to render it.
 */
export declare function $readFile(path: string): string;
