import { Signal } from "./signal";

export class ReadonlySignal<T> extends Signal<T> {
  private static readonly setError = new Error(
    "Cannot set value of a ReadonlySignal"
  );

  // Both of these methods throw an error to prevent changing the value of a
  // ComputedSignal directly.
  // If you want to change the value, you need to change one of its dependencies
  // which will automatically update the ComputedSignal's value.
  //
  // Note that we also use a spread operator to accept any arguments, so that
  // this class doesn't have to be updated if the base Signal class changes.
  public override set(..._args: unknown[]): this {
    throw ReadonlySignal.setError;
  }

  public override update(..._args: unknown[]): this {
    throw ReadonlySignal.setError;
  }
}
