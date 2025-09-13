import { TypeError } from "../conditions/error";

export type Assert<T, Expected> = T extends Expected
  ? true
  : TypeError<"Assertion failed">;
