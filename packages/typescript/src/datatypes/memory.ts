import type { Unwrap } from "../structural/unwrap";
import type { FunctionType } from "./functions";
import type { ObjectType } from "./objects";

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
