/**
 * @description
 * The TypeScript `object` type is really a union between all non-falsy values.
 * This type is a more accurate representation of a plain JavaScript objects
 * that also supports structural typing.
 */
export type ObjectType<Keys extends PropertyKey = string, Values = unknown> = {
  [key in Keys]: Values;
};

/**
 * @description
 * An empty object type that does not allow any properties.
 *
 * This is useful because the `{}` type in TypeScript is special because it
 * allows any non-nullish value.
 */
export type EmptyObject = { [key in never]: never };
