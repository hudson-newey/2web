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

### "When" decorators

"when decorators" document a functions expected calling state, and will
**reject** if the pre-conditions are not reached.

This function will **reject** because it turns the return type into a `Promise`.

```ts
@when([
  assert((person) => iif(() => person.age > 0)),
])
function sayHello(person: Person) {
}
```

## Actions

### "Required" actions

An async callback that MUST be executed within a given timeframe.

If the async queue is in the process of breaking down / crashing / stuck in an
infinite loop, the "required" action will be forcefully moved to the main thread
to prioritize execution.

If the system reaches a critical condition, "optional" actions will be culled in
an attempt to delay system failure.

```ts
async function crashAsyncQueue() {
  setInterval(() => {
    crashAsyncQueue();
  }, 5);
}

crashAsyncQueue();

await new Promise((res) => {
  setTimeout(() => res(), 1_000);
});

// Under normal conditions, this callback would never be invoked, in a realistic
// timeframe.
// However, because it's wrapped in a "required" callback, it will be moved to
// the main thread after 1 second of not being executed.
required(async () => {
  fetch("https://example.com");
});
```

### "Optional" actions

Non-critical functionality that can be culled if the system reaches a "critical"
condition (e.g. overheating, out of memory, etc...)

```ts
optional(async () => {
  console.log(new Date());
});
```

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
