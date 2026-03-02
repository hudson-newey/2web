/**
 * A shorter alias for `Promise<void>`.
 * This allows the emission of any async value, but not the usage or held.
 * This differs from the `unknown` type that can be held.
 */
export type VoidAsync = Promise<void>;
