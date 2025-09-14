export const hiddenType = Symbol("hidden");

/**
 * @description
 * Hides a TypeScripts value type while maintaining a reference to its original
 * type.
 * This is similar to the `never` type, but allows you to move out of the
 * `never` type if it becomes valid again.
 */
export type Hide<T> = { [hiddenType]: T };
