/**
 * Enforces that your bigint datatype cannot accidently use a falsy condition
 * to check if the value is defined.
 */
export type SafeBigInt<T extends bigint> = Exclude<T, 0n> | 0n;
