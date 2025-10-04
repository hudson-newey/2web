/**
 * @description
 * Negates an "If" condition.
 */
export type Not<T, U = true> = T extends U ? never : T;
