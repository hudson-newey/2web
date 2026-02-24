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

  private _value!: Readonly<T>;

  public constructor(initialValue: T) {
    // We use a Promise to defer the initialization to that onCreation and pipes
    // can be set up before the first value is assigned.
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
    return this;
  }

  public update<U extends T>(updater: (currentValue: T) => U): this {
    this.value = updater(this.value);
    return this;
  }

  public subscribe(...callbacks: Subscription<T>[]): this {
    for (const fn of callbacks) {
      this.subscribers.add(fn);
    }
    return this;
  }

  public unsubscribe(...callbacks: Subscription<T>[]): this {
    for (const fn of callbacks) {
      this.subscribers.delete(fn);
    }

    return this;
  }
}

// These void types allow any return type from the callback functions, but will
// enforce that the return type is not used.
type Subscription<T> = (value: T) => void;
