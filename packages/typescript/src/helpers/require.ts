/**
 * @description
 * Makes an optional property in an object required.
 */
export type Require<T, K extends keyof T> = T & {
  [P in K]-?: T[P];
};
