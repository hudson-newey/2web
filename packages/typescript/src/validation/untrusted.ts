import { ObjectType } from "../datatypes/objects";

/**
 * @description
 * Untrusted type that requires you to type narrow it to its correct type
 * before use.
 */
export type Untrusted<ExpectedType = unknown> =
  | ExpectedType
  | string
  | number
  | boolean
  | null
  | undefined
  | symbol
  | bigint
  | ObjectType;
