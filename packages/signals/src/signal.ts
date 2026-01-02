import type { Pipe } from "./pipes/pipe";
import { addExecSubscriber } from "./utils/execCallback";

/**
 * @description
 * Reactive signal that holds a value and notifies subscribers on changes.
 */
export function signal<T>(value: T) {
  return new Signal(value);
}

/**
 * @private
 */
export class Signal<T> {
  private readonly subscribers = new Set<Subscription<T>>();
  private readonly pipes: Pipe<T>[] = [];
  private _value!: Readonly<T>;

  public constructor(initialValue: T) {
    this.value = Object.freeze(initialValue);
  }

  public get value(): T {
    addExecSubscriber(this);

    return this._value;
  }

  protected set value(newValue: T) {
    // To prevent unnecessary updates, we check for strict equality.
    // If the new value is the same as the current value, we do nothing.
    //
    // We use Object.is() instead of === to handle edge cases like NaN and
    // -0/+0 correctly.
    // see: https://stackoverflow.com/a/30543212/10262458
    if (Object.is(this.value, newValue)) {
      return;
    }

    // The _value is frozen to prevent accidental mutations of objects/arrays
    // stored in signals. This ensures that updates to signals are explicit
    // through the set() or update() methods.
    this._value = Object.freeze(newValue);
  }

  public set<U extends T>(newValue: U): this {
    this.value = newValue;
    this.afterChange(newValue);
    return this;
  }

  public update<U extends T>(updater: (currentValue: T) => U): this {
    this.value = updater(this.value);
    return this;
  }

  public subscribe(callback: Subscription<T>): this {
    this.subscribers.add(callback);
    return this;
  }

  public unsubscribe(callback: Subscription<T>): this {
    this.subscribers.delete(callback);
    return this;
  }

  public pipe(...pipeFns: Pipe<T>[]): this {
    this.pipes.push(...pipeFns);
    return this;
  }

  // Lifecycle hook that can be overridden by subclasses
  protected afterChange(newValue: T) {
    for (const subscription of this.subscribers) {
      subscription(newValue);
    }
  }
}

type Subscription<T> = (value: T) => unknown;
