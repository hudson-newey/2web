type Subscription<T> = (value: T) => unknown;

export class Signal<T> {
  private readonly subscribers = new Set<Subscription<T>>();
  private _value!: Readonly<T>;

  public constructor(initialValue: T) {
    this.value = Object.freeze(initialValue);
  }

  public get value(): T {
    return this._value;
  }

  protected set value(newValue: T) {
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

  public subscribe(callback: Subscription<T>): this {
    this.subscribers.add(callback);
    return this;
  }

  public unsubscribe(callback: Subscription<T>): this {
    this.subscribers.delete(callback);
    return this;
  }

  // Lifecycle hook that can be overridden by subclasses
  protected afterChange(newValue: T) {
    for (const subscription of this.subscribers) {
      subscription(newValue);
    }
  }
}
