import { ReadonlySignal } from "../readonlySignal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal that hooks into a ResizeObserver and updates whenever the size
 * of the queried element changes.
 */
export function resizeObserver<ObservedElements extends MaybeSignal<Element[]>>(
  observedElements: ObservedElements,
) {
  return new ResizeObserverSignal().init(observedElements);
}

class ResizeObserverSignal<
  ObservedElements extends MaybeSignal<Element[]>,
> extends ReadonlySignal<ResizeObserverEntry[]> {
  public constructor() {
    super([]);
  }

  public async init(observedElements: ObservedElements) {
    const observedTargets = new Set<Element>;

    const observer = new ResizeObserver(
      async (entries: ResizeObserverEntry[]) => {
        this._internalSet(entries);
      },
    );

    const observe = (element: Element) => {
      if (observedTargets.has(element)) return;

      observedTargets.add(element);
      observer.observe(element);
    };

    const targets = await unwrapSignal(observedElements);
    for (const element of targets) {
      observe(element);
    }

    if (isSignal(observedElements)) {
      observedElements.subscribe((elements) => {
        if (!elements) return;

        for (const element of elements) {
          observe(element);
        }
      });
    }
  }
}
