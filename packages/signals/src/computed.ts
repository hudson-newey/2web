import { ReadonlySignal } from "./readonlySignal";
import { Signal } from "./signal";

export type ComputedSignalReducer<T> = () => T;

export class ComputedSignal<T> extends ReadonlySignal<T> {
  public constructor(
    reducer: ComputedSignalReducer<T>,
    dependencies: Signal<unknown>[],
  ) {
    super(reducer());

    for (const dependency of dependencies) {
      dependency.subscribe(() => {
        this.value = reducer();
      });
    }
  }
}
