import { describeType, expectType } from "./expect";

describeType("true", () => {
  expectType(true).toBeTrue();

  // Expect failures
  expectType(false).toBeTrue();
  expectType(1).toBeTrue();
});

describeType("false", () => {
  expectType(false).toBeFalse();

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

  // Expect failures
  expectType("hello").toBeNumber();
  expectType(true).toBeNumber();
});

describeType("string", () => {
  expectType("hello").toBeString();
  expectType("").toBeString();

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

  // Expect failures
  expectType(null).toBeUndefined();
  expectType(0).toBeUndefined();
});

describeType("functions", () => {
  expectType(() => {}).toBeFunction();
  expectType(function () {}).toBeFunction();
  expectType(() => 42).toBeFunction();
  expectType(async () => {}).toBeFunction();
  expectType(function* () {}).toBeFunction();
  expectType({}.toString).toBeFunction();

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
  expectType(new Array()).toBeArray();

  // Expect failures
  expectType({}).toBeArray();
  expectType("hello").toBeArray();
  expectType(true).toBeArray();
  expectType(new Set()).toBeArray();
});

describeType("regexps", () => {
  expectType(/abc/).toBeRegExp();
  expectType(new RegExp("abc")).toBeRegExp();

  // Expect failures
  expectType({}).toBeRegExp();
  expectType("abc").toBeRegExp();
  expectType(12).toBeRegExp();
  expectType([]).toBeRegExp();
});
