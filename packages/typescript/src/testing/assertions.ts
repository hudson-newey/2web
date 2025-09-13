import type { CompileError } from "../conditions/error";

export type Assert<T, Expected> = T extends Expected
  ? true
  : CompileError<"Assertion failed">;
