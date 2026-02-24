import type { OnCreateCallback } from "./lifecycle/onBeforeCreate";
import type { OnFirstUpdateCallback } from "./lifecycle/onFirstUpdate";
import type { OnDestroyCallback } from "./lifecycle/onDestroy";
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
  // This value will be used before the signals initial update
  public static readonly unsetValue = null;

  private readonly subscribers = new Set<Subscription<typeof this.value>>();

  // Lifecycle and pipe callbacks are an array instead of a Set to preserve
  // order of execution and so that duplicates are allowed.
  private readonly onCreateCallbacks: OnCreateCallback[] = [];
  private readonly onFirstUpdateCallbacks: OnFirstUpdateCallback<T>[] = [];
  private readonly onDestroyCallbacks: OnDestroyCallback<T>[] = [];
  private readonly pipes: Pipe<T>[] = [];

  // Because the signal ctor sets the initial value in the async queue, it is
  // possible for the signals value to be null before the signals initial value
  // is set.
  private _value: Readonly<T> | typeof Signal.unsetValue = Signal.unsetValue;

  // This boolean is used to track if we have already run the onCreate
  // callbacks.
  // I use this to issue a warning if the onCreate callbacks are added after
  // the signal has already been created.
  private doneFirstUpdate = false;
  private doneCreate = false;

  public constructor(initialValue: typeof this.value = Signal.unsetValue) {
    // Before we set any value, we run the beforeCreate callbacks
    for (const beforeCreateFn of this.onCreateCallbacks) {
      beforeCreateFn();
    }

    this.doneCreate = true;

    // We use a Promise to defer the initialization to that onCreation and pipes
    // can be set up before the first value is assigned.
    this._internalSet(initialValue).then(() => {
      if (initialValue !== Signal.unsetValue) {
        // We set the doneFirstUpdate before running the onCreate callbacks so if
        // there is a hard fail in one of the callbacks, the doneFirstUpdate flag
        // is still set.
        this.doneFirstUpdate = true;

        for (const firstUpdateFn of this.onFirstUpdateCallbacks) {
          firstUpdateFn(this.value);
        }
      }
    });
  }

  public get value(): T | null {
    addExecSubscriber(this);

    return this._value;
  }

  public set<U extends T>(newValue: U): Promise<this> {
    return this._internalSet(newValue);
  }

  public update<U extends T>(updater: UpdaterFn<T, U>): Promise<this> {
    const newValue = updater(this.value);
    return this.set(newValue);
  }

  public subscribe(...callbacks: Subscription<typeof this.value>[]): this {
    for (const fn of callbacks) {
      this.subscribers.add(fn);
    }

    return this;
  }

  public unsubscribe(...callbacks: Subscription<typeof this.value>[]): this {
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
  public onCreate(...callbacks: OnCreateCallback[]): this {
    if (this.doneCreate) {
      console.warn(
        "[@two-web/kit/signals] onCreate callbacks added after the signal " +
          "has already completed the onCreate lifecycle. These callbacks " +
          "will never be executed.",
      );
    }

    this.onCreateCallbacks.push(...callbacks);
    return this;
  }

  /**
   * Executed only once, the first time the signal's value is accessed.
   * The signal is returned to allow method chaining.
   */
  public onFirstUpdate(...callbacks: OnFirstUpdateCallback<T>[]): this {
    if (this.doneFirstUpdate) {
      console.warn(
        "[@two-web/kit/signals] onFirstUpdate callbacks added after the " +
          "signal has already completed the onFirstUpdate lifecycle. " +
          "These callbacks will never be executed.",
      );
    }

    this.onFirstUpdateCallbacks.push(...callbacks);
    return this;
  }

  /**
   * Executed only once, when the signal is destroyed.
   * The signal is returned to allow method chaining.
   * You will need to use the "using" keyword to hook into the signal's
   * onDestroy lifecycle.
   */
  public onDestroy(...callbacks: OnDestroyCallback<T>[]): this {
    this.onDestroyCallbacks.push(...callbacks);
    return this;
  }

  protected _internalSet(newValue: typeof this.value): Promise<this> {
    // To prevent unnecessary updates, we check for strict equality.
    // If the new value is the same as the current value, we do nothing.
    //
    // We use Object.is() instead of === to handle edge cases like NaN and
    // -0/+0 correctly.
    // see: https://stackoverflow.com/a/30543212/10262458
    if (Object.is(this.value, newValue)) {
      return Promise.resolve(this);
    }

    return new Promise((res) => {
      // Use a setTimeout here so that you can chain multiple listeners to an
      // update before setting the value.
      setTimeout(() => {
        if (!this.doneFirstUpdate) {
          this.doneFirstUpdate = true;
          for (const firstUpdateFn of this.onFirstUpdateCallbacks) {
            firstUpdateFn(this.value);
          }
        }

        // The _value is frozen to prevent accidental mutations of objects/arrays
        // stored in signals. This ensures that updates to signals are explicit
        // through the set() or update() methods.
        this._value = Object.freeze(newValue);
        this._afterChange(newValue);

        res(this);
      }, 0);
    });
  }

  // Lifecycle hook that can be overridden by subclasses
  protected _afterChange(newValue: typeof this.value): void {
    for (const subscription of this.subscribers) {
      subscription(newValue);
    }
  }

  // This is a special method that is called when the signal goes out of scope
  // when the "using" keyword is used.
  // It is used to run the onDestroy callbacks.
  protected [Symbol.dispose]() {
    for (const callback of this.onDestroyCallbacks) {
      callback(this.value);
    }
  }
}

// This void type allows any return type from the callback functions, but will
// enforce that the return type is not used.
export type Subscription<T> = (value: T) => void;
export type UpdaterFn<T, U extends T> = (currentValue: T | null) => U;
