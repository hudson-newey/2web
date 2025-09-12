import { Unwrap } from "../structural/unwrap";
import { FunctionType } from "./functions";
import { ObjectType } from "./objects";

export type StackVariable = string | number | boolean | null | undefined | symbol | bigint;
export type HeapVariable = unknown[] | FunctionType | ObjectType;
export type Variable = Unwrap<StackVariable | HeapVariable>;

export type StringTemplatable = string | number | bigint | boolean | null | undefined
