import type { Unwrap } from "../structural/unwrap.ts";
import type { FunctionType } from "./functions.ts";
import type { ObjectType } from "./objects.ts";

export type StackVariable =
  | string
  | number
  | boolean
  | null
  | undefined
  | symbol
  | bigint;
export type HeapVariable = unknown[] | FunctionType | ObjectType;
export type Variable = Unwrap<StackVariable | HeapVariable>;

export type StringTemplatable =
  | string
  | number
  | bigint
  | boolean
  | null
  | undefined;
