# 2Web Kit - Iron Guard

Defensive programming helpers.

If you're using TypeScript, you probably don't need this library unless you're
creating high-reliability software that must never crash.

## Errors

### Retry

Retries a statement until it passes

```ts
// Attempt to fetch the current user model 5 times, if one request fails, retry
// after 1 second.
// If 5 attempts fail, the last error will be thrown.
const user = await retry(() => {
  return await fetch("https://example.com");
}, { interval: 1_000, times: 5 });
```

### Wait

An awaitable function to unwrap promises while handling their rejections.

This is needed because a normal async function does have a `catch` callback, but
this cannot be easily used with the `await` keyword.

```ts
const resp = await wait(fetch("https://example.com"), {
  catch: (err) => {
    throw new Error("Error fetching user data", { context: err });
  }
});
```

### Try Catch

Requires you to handle error cases.

```ts
const resp = tryCatch(() => throw Error("my error"));
if (resp instanceof Error) {
  console.log("handle the error");
}
```

## Conditionals

When creating high-reliability software, type guards can cause the program to
hard fail if an unexpected type is reached.

```ts
const x = {};

// The following would throw an error with a normal "if" condition (because
// "person" is undefined), causing the program to hard-fail.
// If you do not want the program to hard fail, you probably want iif.
//
// If your conditional callback is async, you can await the iif call.
iif(() => x.person.age > 0)({
  then: () => {},
  else: () => {},
  catch: () => {},
})
```

## Assert

Asserts that the program is in the expected state before continuing.

This helper is intended to document the expected state of the program, not catch
errors or be used as a guard or in control flow.

```ts
function () {
  assert(() => 5 > 2);
}
```

## Actions

### Check

Performs an action twice and asserts that the results are the same.

```ts
// Ensure that the operation inside the callback is something that JIT cannot
// optimize away by re-using the same value.
// Additionally, this callback should have no side-effects.
check(() => {
  fetch("https://example.com");
});
```

## Ephemeral

Creates an ephemeral stack frame to that is automatically zeroed out after use.

This is useful for handling sensitive data such as passwords or cryptographic
keys where you don't want the data to linger in memory longer than necessary.

```ts
// After this scope ends, the secret variable will be zeroed out.
{
  const secret = using ephemeral("my super secret password");
  // Use the secret here
}
```
