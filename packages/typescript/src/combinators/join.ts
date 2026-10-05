import type { CompileError } from "../conditions/error.ts";
import type { Extends } from "../conditions/extends.ts";
import type { Or } from "../conditions/or.ts";

/**
 * @description
 * A self-documenting polymorphic type.
 */
export type Join<A, B> = Or<Extends<A, B>, Extends<B, A>> extends true
  ? B & A
  : CompileError<`Types are not compatible`>;
