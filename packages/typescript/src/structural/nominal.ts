/**
 * @description
 * Converts a structural type to a nominal type
 */
export type Nominal<T> = T extends infer U ? U : never;
