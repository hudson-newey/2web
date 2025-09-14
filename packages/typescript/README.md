# 2Web Kit - TypeScript Helpers

A large set of utility types for making high-quality, resilient applications
without having to maintain bespoke type definitions.

## Assertions

This package contains methods to assert TypeScript types against real javascript
values.

```ts
import { expectType } from "@two-web/kit/packages/typescript";

type Person = { name: string, age: number };
type Shape = { width: number, height: number, depth: number };

expectType(true).toBeTrue();
expectType(true).not.toBeFalse();

expectType({ name: "bob", age: 52 }).toBe<Person>();
expectType({ name: "bob", age: 52 }).not.toBe<Shape>();
```

## Untrusted Types

### `Untrusted<T>`

User input cannot be trusted, and therefore requires type narrowing for robust
applications.

For this case, you can use the `Untrusted<T>` type which requires you to type
narrow before usage.

### `Unstable<T>`

When creating library or extension code, we cannot reliably know what the host
application / browser is going to do.
It might remove/clone element nodes that we previously created, so we can't
reliably keep element or object references.

In these cases, you do not want to type narrow outside of usage, which is
exactly what the `Unstable<T>` type does.

You can use the `open` function to use the unstable type, but you will need
to check that your assumptions hold before every usage because the type
narrowing is not leaked outside of the `open` callback.

```ts
// The "element" type is "never" to ensure that we cannot use it outside of a
// safe "open" callback.
const element = unstable(document.querySelector(".my-element"));

// We are forced to use the "open" function to type narrow the unstable
// type to a stable type.
open(element, (stableElement) => {
  if (!stableElement) {
    throw new Error("Element not found");
  }

  stableElement.classList.add("active");
});

// Note that even though the "element" was narrowed inside of the "open"
// callback using a type guard, it is still "never" out here, because the host
// application/browser/extension might have removed the element from the DOM.
```

## Type Helpers

TypeScript has a lot of quirks such as the `object` type really being for any
type that is not `null` or `undefined`.

There are some type helpers in the two-web kit package to make these definitions
less verbose.

```ts
import { FunctionType, ObjectType } from "@two-web/kit/packages/typescript";

type Person = ObjectType<string, unknown>;
type SearchPerson = FunctionType<[id: number], Person>;
```
