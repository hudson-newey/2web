/**
 * Enforces that your string datatype cannot accidently use a falsy condition
 * to check if the value is defined.
 */
export type SafeString<T extends string> = Omit<T, ""> | "";
