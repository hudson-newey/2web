import type { CompileError } from "../conditions/error.ts";

export type Assert<T, Expected> = T extends Expected
  ? true
  : CompileError<"Assertion failed">;
