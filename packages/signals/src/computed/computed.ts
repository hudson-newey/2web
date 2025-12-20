import { ReadonlySignal } from "../readonlySignal";
import type { Signal } from "../signal";

export function computed<T>(
  reducer: ComputedSignalReducer<T>,
  dependencies: Signal<unknown>[],
) {
  return new ComputedSignal(reducer, dependencies);
}

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

type ComputedSignalReducer<T> = () => T;
