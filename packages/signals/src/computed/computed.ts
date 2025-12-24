import { ReadonlySignal } from "../readonlySignal";
import type { Signal } from "../signal";
import { execCallback } from "../utils/execCallback";

export function computed<T>(reducer: ComputedSignalReducer<T>) {
  const { subscribers } = execCallback(reducer);

  return new ComputedSignal(reducer, subscribers);
}

class ComputedSignal<T> extends ReadonlySignal<T> {
  public constructor(
    reducer: ComputedSignalReducer<T>,
    dependencies: ReadonlySet<Signal<unknown>>
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
