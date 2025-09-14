import type { Without } from "../helpers/without";

/**
 * @description
 * Autocompletes a type `T` with a set of suggested types `Suggestions` without
 * narrowing the type to just the suggestions.
 *
 */
export type Autocomplete<T, Suggestions extends T> =
  | Suggestions
  | Without<T, Suggestions>;
