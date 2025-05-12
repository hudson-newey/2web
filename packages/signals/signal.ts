type Subscription<T> = (value: T) => unknown;

export class Signal<T> {
  public constructor(initialValue: T) {
    this._value = initialValue;
  }

  private subscriptions = new Set<Subscription<T>>();
  private _value: T;

  public get value(): T {
    return this._value;
  }

  public set(newValue: T) {
    this._value = newValue;
    this.update();
  }

  public subscribe(callback: Subscription<T>) {
    this.subscriptions.add(callback);
  }

  public unsubscribe(callback: Subscription<T>) {
    this.subscriptions.delete(callback);
  }

  private update() {
    for (const subscription of this.subscriptions) {
      subscription(this.value);
    }
  }
}
