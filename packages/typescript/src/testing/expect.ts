import { ObjectType } from "../datatypes/objects";
import { Assert } from "./assertions";
import { FunctionType } from "../datatypes/functions";
import { ErrorType } from "../conditions/error";

type AssertMethod<T, Expected> = Assert<T, Expected> extends true
  ? FunctionType<[], true>
  : ErrorType<`Assertion failed.`>;

interface AssertionMethods<T> {
  toBeBoolean: AssertMethod<T, boolean>;
  toBeTrue: AssertMethod<T, true>;
  toBeFalse: AssertMethod<T, false>;

  toBeNumber: AssertMethod<T, number>;
  toBeString: AssertMethod<T, string>;
  toBeBigInt: AssertMethod<T, bigint>;

  toBeSymbol: AssertMethod<T, symbol>;
  toBeUndefined: AssertMethod<T, undefined>;
  toBeNull: AssertMethod<T, null>;
  toBeFunction: AssertMethod<T, Function>;
  toBeObject: AssertMethod<T, ObjectType>;
  toBeArray: AssertMethod<T, unknown[]>;
  toBeRegExp: AssertMethod<T, RegExp>;
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
 */
export const expectType = <const T>(_value: T): AssertionMethods<T> => {
  return {
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
  };
};

export const describeType = <const Description extends string>(
  _description: Description,
  _fn: () => void,
) => {};
