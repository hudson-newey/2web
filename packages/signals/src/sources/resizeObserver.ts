import { ReadonlySignal } from "../readonlySignal";

/**
 * @description
 * A signal that hooks into a ResizeObserver and updates whenever the size
 * of the queried element changes.
 */
export class ResizeObserverSignal<
  ObservedElements extends Element[]
> extends ReadonlySignal<ResizeObserverEntry[]> {
  public constructor(observedElements: ObservedElements) {
    super([]);

    const callback: ResizeObserverCallback = (entries: ResizeObserverEntry[]) => {
      this.set(entries);
    };

    const observer = new ResizeObserver(callback);
    for (const element of observedElements) {
      observer.observe(element);
    }
  }
}
