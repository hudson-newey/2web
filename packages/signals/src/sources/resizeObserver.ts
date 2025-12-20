import { ReadonlySignal } from "../readonlySignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal that hooks into a ResizeObserver and updates whenever the size
 * of the queried element changes.
 */
export function resizeObserver<ObservedElements extends MaybeSignal<Element[]>>(
  observedElements: ObservedElements
) {
  return new ResizeObserverSignal(observedElements);
}

class ResizeObserverSignal<
  ObservedElements extends MaybeSignal<Element[]>
> extends ReadonlySignal<ResizeObserverEntry[]> {
  public constructor(observedElements: ObservedElements) {
    super([]);

    const callback: ResizeObserverCallback = (
      entries: ResizeObserverEntry[]
    ) => {
      this.value = entries;
    };

    const targets = unwrapSignal(observedElements);
    const observer = new ResizeObserver(callback);
    for (const element of targets) {
      observer.observe(element);
    }
  }
}
