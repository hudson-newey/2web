import { expect, test } from "vitest";
import { Signal } from "./signal";

test("should be able to change a signals value", () => {
  const mySignal = new Signal(0);

  mySignal.set(1);
  expect(mySignal.value).toBe(1);

  // Test updating through the "update" method
  mySignal.update((v) => v + 1);
  expect(mySignal.value).toBe(2);
});

test("should notify subscribers when the value changes", () => {
  const mySignal = new Signal("abc");

  let notifiedValue: string | null = null;
  const subscription = (value: string) => {
    notifiedValue = value;
  };

  mySignal.subscribe(subscription);
  mySignal.set("def");

  expect(notifiedValue).toBe("def");

  // Test that unsubscribe works
  mySignal.unsubscribe(subscription);
  mySignal.set("ghi");
  expect(notifiedValue).toBe("def"); // Should not have changed
});
