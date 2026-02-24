import { expect, test } from "vitest";
import { signal } from "./signal";
import { onFirstUpdate } from "./lifecycle/onFirstUpdate";

test("should be able to change a signals value", async () => {
  const mySignal = signal(0);
  await onFirstUpdate(mySignal);

  await mySignal.set(1);
  expect(mySignal.value).toBe(1);

  // Test updating through the "update" method
  await mySignal.update((v) => v + 1);
  expect(mySignal.value).toBe(2);
});

test("should notify subscribers when the value changes", async () => {
  const mySignal = signal("abc");
  await onFirstUpdate(mySignal);

  let notifiedValue: string | null = null;
  const subscription = (value: string) => {
    notifiedValue = value;
  };

  mySignal.subscribe(subscription);
  await mySignal.set("def");

  expect(notifiedValue).toBe("def");

  // Test that unsubscribe works
  mySignal.unsubscribe(subscription);
  await mySignal.set("ghi");
  expect(notifiedValue).toBe("def"); // Should not have changed
});
