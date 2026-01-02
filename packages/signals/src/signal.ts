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
  private readonly beforeCreateCallbacks = new Set<BeforeCreateCallback>();
  private readonly onCreateCallbacks = new Set<OnCreateCallback<T>>();
  private readonly onDestroyCallbacks = new Set<OnDestroyCallback<T>>();

  // Pipes are an array instead of a Set to preserve order of execution and so
  // that duplicates are allowed.
  private readonly pipes: Pipe<T>[] = [];
  private _value!: Readonly<T>;

  public constructor(initialValue: T) {
    // Before we set any value, we run the beforeCreate callbacks
    for (const callback of this.beforeCreateCallbacks) {
      callback();
    }

    // We use a Promise to defer the initialization to that onCreation and pipes
    // can be set up before the first value is assigned.
    new Promise(() => {
      this.value = Object.freeze(initialValue);

      // Run onCreate callbacks
      for (const callback of this.onCreateCallbacks) {
        callback(this.value);
      }
    });
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
    this.afterChange(newValue);
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

  public pipe(...pipeFns: Pipe<T>[]): this {
    this.pipes.push(...pipeFns);
    return this;
  }

  /**
   * Runs before the signal is created and before any value is set.
   */
  public beforeCreate(...callbacks: BeforeCreateCallback<T>[]): this {
    for (const fn of callbacks) {
      this.onCreateCallbacks.add(fn);
    }

    return this;
  }

  /**
   * Executed only once, the first time the signal's value is accessed.
   * The signal is returned to allow method chaining.
   */
  public onCreate(...callbacks: OnCreateCallback<T>[]): this {
    for (const fn of callbacks) {
      this.onCreateCallbacks.add(fn);
    }

    return this;
  }

  /**
   * Executed only once, when the signal is destroyed.
   * The signal is returned to allow method chaining.
   * You will need to use the "using" keyword to hook into the signal's
   * onDestroy lifecycle.
   */
  public onDestroy(...callbacks: OnDestroyCallback<T>[]): this {
    for (const fn of callbacks) {
      this.subscribers.delete(fn);
    }

    return this;
  }

  // Lifecycle hook that can be overridden by subclasses
  protected afterChange(newValue: T): void {
    for (const subscription of this.subscribers) {
      subscription(newValue);
    }
  }

  protected [Symbol.dispose]() {
    for (const callback of this.onDestroyCallbacks) {
      callback(this.value);
    }
  }
}

// These void types allow any return type from the callback functions, but will
// enforce that the return type is not used.
type Subscription<T> = (value: T) => void;
type BeforeCreateCallback = () => void;
type OnCreateCallback<T> = (value: T) => void;
type OnDestroyCallback<T> = (value: T) => void;
