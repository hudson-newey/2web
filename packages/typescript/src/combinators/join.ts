import type { CompileError } from "../conditions/error";
import type { Extends } from "../conditions/extends";
import type { Or } from "../conditions/or";

/**
 * @description
 * A self-documenting polymorphic type.
 */
export type Join<A, B> = Or<Extends<A, B>, Extends<B, A>> extends true
  ? B & A
  : CompileError<`Types are not compatible`>;
