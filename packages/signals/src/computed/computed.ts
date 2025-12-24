import { ReadonlySignal } from "../readonlySignal";
import { execCallback } from "../utils/execCallback";

export function computed<T>(reducer: ComputedSignalReducer<T>) {
  return new ComputedSignal(reducer);
}

class ComputedSignal<T> extends ReadonlySignal<T> {
  public constructor(reducer: ComputedSignalReducer<T>) {
    const { dependencies } = execCallback(reducer);

    super(reducer());

    for (const dependency of dependencies) {
      dependency.subscribe(() => {
        this.value = reducer();
      });
    }
  }
}

type ComputedSignalReducer<T> = () => T;
