/**
 * @description
 * The default "Function" type in TypeScript has some quirks, and it's therefore
 * advised to use the anonymous function type instead.
 * However, this is very verbose, leads to duplicating type information, and
 * prone to a development burden if structural typing is involved.
 */
export type FunctionType<Args extends unknown[] = unknown[], ReturnType = unknown> = (
  ...args: Args
) => ReturnType;
