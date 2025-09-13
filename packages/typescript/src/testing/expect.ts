import { ObjectType } from "../datatypes/objects";
import { Assert } from "./assertions";
import { FunctionType } from "../datatypes/functions";
import { TypeError } from "../conditions/error";

type AssertMethod<T, Expected, ExpectPass = true> = Assert<
  T,
  Expected
> extends ExpectPass
  ? FunctionType<[], true>
  : TypeError<`Assertion failed`, { expected: Expected, found: T }>;

interface Modifiers<T> {
  not: AssertionMethods<T, TypeError<"Assertion failed">>;
}

interface AssertionMethods<T, ExpectPass = true>
  extends Modifiers<T> {
  toBeBoolean: AssertMethod<T, boolean, ExpectPass>;
  toBeTrue: AssertMethod<T, true, ExpectPass>;
  toBeFalse: AssertMethod<T, false, ExpectPass>;

  toBeNumber: AssertMethod<T, number, ExpectPass>;
  toBeString: AssertMethod<T, string, ExpectPass>;
  toBeBigInt: AssertMethod<T, bigint, ExpectPass>;

  toBeSymbol: AssertMethod<T, symbol, ExpectPass>;
  toBeUndefined: AssertMethod<T, undefined, ExpectPass>;
  toBeNull: AssertMethod<T, null, ExpectPass>;
  toBeFunction: AssertMethod<T, Function, ExpectPass>;
  toBeObject: AssertMethod<T, ObjectType, ExpectPass>;
  toBeArray: AssertMethod<T, unknown[] | readonly unknown[], ExpectPass>;
  toBeRegExp: AssertMethod<T, RegExp, ExpectPass>;
}

const incorrectEnvironmentError = (() =>
  new Error(
    "This method is only for type assertions and should not be called at runtime",
  )) as never;

/**
 * @description
 * A compile-time assertion function that can be used to write TypeScript type
 * tests.
 *
 * @example
 * ```ts
 * expectType(true).toBeTrue();
 * expectType(false).toBeTrue();
 * // ^ Will cause a type error
 *
 * expectType(12).toBeNumber();
 * expectType("hello").toBeString();
 *
 * expectType(() => 42).toReturn(42);
 * expectType(() => 42).toReturn(Number);
 * expectType(() => 42).toReturn<Number>();
 * ```
 *
 * We provide a bail-out for custom types that are not covered by the built-in
 * assertions.
 *
 * ```ts
 * type User = { id: number; name: string };
 *
 * const myUser: User = { id: 1, name: "Alice" };
 * const myLocation = { lat: 20.0, lng: 20.25 };
 *
 * expectType(myUser).toBe<User>();
 * expectType(myUser).toBe<User>(); // This will cause a type error
 * ```
 *
 * There are also "negative" assertions which are useful if you want a type to
 * fail under certain conditions.
 *
 * ```ts
 * type User = { id: number; name: string };
 * type AdminUser = User & { isAdmin: true };
 *
 * const myUser: User = { id: 1, name: "Alice" };
 *
 * expect(myUser).not.toBe<AdminUser>();
 * // This will not cause a type error because we expect the type to fail
 * ```
 *
 * ## Type Guards
 *
 * You can also use the `expectType` function as a type guard to assert that a
 * variable is a certain type at runtime.
 *
 * ```ts
 * const value = true;
 * if (expectType(value).toBeTrue()) {
 *   console.log("The value is true");
 * }
 *
 * // TypeScript can automatically determine that the following code block is
 * // unreachable because `false` is not `true`.
 * if (expectType(value).toBeFalse()) {
 *   // This code is unreachable
 * }
 * ```
 */
export const expectType = <const T>(_value?: T): AssertionMethods<T> => {
  const assertions = {
    toBeBoolean: incorrectEnvironmentError,
    toBeTrue: incorrectEnvironmentError,
    toBeFalse: incorrectEnvironmentError,

    toBeNumber: incorrectEnvironmentError,
    toBeString: incorrectEnvironmentError,
    toBeBigInt: incorrectEnvironmentError,

    toBeSymbol: incorrectEnvironmentError,
    toBeUndefined: incorrectEnvironmentError,
    toBeNull: incorrectEnvironmentError,
    toBeFunction: incorrectEnvironmentError,
    toBeObject: incorrectEnvironmentError,
    toBeArray: incorrectEnvironmentError,
    toBeRegExp: incorrectEnvironmentError,
  } as any;

  assertions.not = assertions;

  return assertions;
};

export const describeType = <const Description extends string>(
  _description: Description,
  _fn: () => void,
) => {};
