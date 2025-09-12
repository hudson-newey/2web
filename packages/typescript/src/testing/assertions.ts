import { ErrorType } from "../conditions/error";

export type Assert<T, Expected> = T extends Expected
  ? true
  : ErrorType<"Assertion failed">;
