// @ts-nocheck

import { describeType, expectType } from "./expect";

describeType("true", () => {
  expectType(true).toBeTrue();

  // Not conditions
  expectType(false).not.toBeTrue();
  expectType(1).not.toBeTrue();
  expectType("hello").not.toBeTrue();

  // Expect failures
  expectType(false).toBeTrue();
  expectType(1).toBeTrue();
});

describeType("false", () => {
  expectType(false).toBeFalse();

  // Not conditions
  expectType(true).not.toBeFalse();
  expectType(0).not.toBeFalse();
  expectType("").not.toBeFalse();

  // Expect failures
  expectType(true).toBeFalse();
  expectType(0).toBeFalse();
  expectType("").toBeFalse();
});

describeType("number", () => {
  expectType(12).toBeNumber();

  // TypeScript has some quirks such as `NaN` and `Infinity` being of type
  // `number`.
  expectType(Infinity).toBeNumber();
  expectType(NaN).toBeNumber();

  // Not conditions
  expectType("24").not.toBeNumber();
  expectType(true).not.toBeNumber();

  // Expect failures
  expectType("hello").toBeNumber();
  expectType(true).toBeNumber();
});

describeType("string", () => {
  expectType("hello").toBeString();
  expectType("").toBeString();
  expectType(String()).toBeString();

  // Not conditions
  expectType(12).not.toBeString();
  expectType(true).not.toBeString();
  expectType({}).not.toBeString();
  expectType(String).not.toBeString();

  // Expect failures
  expectType(12).toBeString();
  expectType(true).toBeString();
});

describeType("null", () => {
  expectType(null).toBeNull();

  // Expect failures
  expectType(undefined).toBeNull();
  expectType(0).toBeNull();
});

describeType("undefined", () => {
  expectType(undefined).toBeUndefined();

  // Not conditions
  expectType(null).not.toBeUndefined();
  expectType(0).not.toBeUndefined();
  expectType("").not.toBeUndefined();

  // Expect failures
  expectType(null).toBeUndefined();
  expectType(0).toBeUndefined();
});

describeType("functions", () => {
  expectType(() => {}).toBeFunction();
  expectType(() => {}).toBeFunction();
  expectType(() => 42).toBeFunction();
  expectType(async () => {}).toBeFunction();
  expectType(function* () {}).toBeFunction();
  expectType({}.toString).toBeFunction();
  expectType(String).toBeFunction();

  // Not conditions
  expectType(12).not.toBeFunction();
  expectType("hello").not.toBeFunction();
  expectType(true).not.toBeFunction();
  expectType({}).not.toBeFunction();
  expectType([]).not.toBeFunction();
  expectType(null).not.toBeFunction();
  expectType(undefined).not.toBeFunction();
  expectType(/abc/).not.toBeFunction();
  expectType(new Date()).not.toBeFunction();

  // Expect failures
  expectType(12).toBeFunction();
  expectType("hello").toBeFunction();
  expectType(true).toBeFunction();
  expectType({}).toBeFunction();
});

describeType("objects", () => {
  expectType({ a: 1, b: "hello" }).toBeObject();
  expectType({}).toBeObject();
  expectType(new Date()).toBeObject();
  expectType(/abc/).toBeObject();
  expectType([]).toBeObject();
  expectType(() => {}).toBeObject();

  // Not conditions
  expectType(null).not.toBeObject();
  expectType(undefined).not.toBeObject();
  expectType(12).not.toBeObject();
  expectType("hello").not.toBeObject();
  expectType(true).not.toBeObject();
  expectType(Symbol()).not.toBeObject();
  expectType(42n).not.toBeObject();

  // Expect failures
  expectType(null).toBeObject();
  expectType(undefined).toBeObject();
  expectType(12).toBeObject();
  expectType("hello").toBeObject();
  expectType(true).toBeObject();
});

describeType("arrays", () => {
  expectType([]).toBeArray();
  expectType([1, 2, 3]).toBeArray();
  expectType([]).toBeArray();

  // Expect failures
  expectType({}).toBeArray();
  expectType("hello").toBeArray();
  expectType(true).toBeArray();
  expectType(new Set()).toBeArray();
});

describeType("regexps", () => {
  expectType(/abc/).toBeRegExp();
  expectType(/abc/).toBeRegExp();

  // Expect failures
  expectType({}).toBeRegExp();
  expectType("abc").toBeRegExp();
  expectType(12).toBeRegExp();
  expectType([]).toBeRegExp();
});
