/**
 * @description
 * The TypeScript `object` type is really a union between all non-falsy values.
 * This type is a more accurate representation of a plain JavaScript objects
 * that also supports structural typing.
 */
export type ObjectType<Keys extends PropertyKey = string, Values = any> = {
  [key in Keys]: Values;
};
