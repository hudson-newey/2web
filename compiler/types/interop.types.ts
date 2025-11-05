// These type declarations allow for strict type checking and interop with the
// 2web compiler.
// Because the compiler exposes a "virtual" imports and keywords that don't
// actually have a runtime representation, we declare them here so that they
// only exist at compile time.
//
// We use a .ts file instead of a .d.ts file so that we get full type checking.
// So that tsconfig.json files can import this file, we re-export everything in
// this file from a d.ts file of the same name.

// This module declares the types for the "@two-web/compiler" virtual module.
// You can import these compiler functions and use these types in your code,
// but they will be stripped out at compile time.
// E.g. import { $uid } from "@two-web/compiler";
declare module "@two-web/compiler" {
  /**
   * @description
   * Generates a unique identifier that is guaranteed to be unique within the
   * current page.
   */
  export function $uid(): string;

  /**
   * @description
   * An object of all of the props passed to the component either as attributes
   * or via property binding.
   *
   * @example
   * ```html
   * <!-- index.2web -->
   * <my-message title="Hello" count="5"></my-message>
   * ```
   *
   * ```ts
   * // message.component.2web
   * <script compiled>
   *   import { $props } from "@two-web/compiler";
   *   $ count = $props().count;
   * </script>
   *
   * <h1>{{ $props().title }}</h1>
   * <p>You have {{ $count }} new messages.</p>
   * ```
   */
  export function $props(): Readonly<Record<string, Readonly<unknown>>>;

  /**
   * @description
   * Retrieves the value of an environment variable at compile time.
   * Returns `undefined` if the environment variable is not defined.
   *
   * @example
   * ```ts
   * import { $env } from "@two-web/compiler";
   *
   * const apiUrl = $env("API_URL") ?? "https://default.api.com";
   * ```
   */
  export function $env(key: string): string | undefined;
}

declare global {
  /**
   * @description
   * A reactive compile-time variable.
   *
   * @example
   * ```html
   * <span #counter-output>{{ $count }}</span>
   *
   * <script compiled>
   *   $ count = 0;
   * </script>
   *
   * <button @click="$count = $count + 1">Increment</button>
   * <button @click="$count = $count - 1">Decrement</button>
   * ```
   */
  const $: any;
}

export {};
