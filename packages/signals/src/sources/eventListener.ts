import { ReadonlySignal } from "../readonlySignal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export function eventListener<
  EventName extends keyof ElementEventMap,
  Target extends MaybeSignal<Element>,
  Callback extends EventListenerCallback,
>(eventName: EventName, element: Target, callback: Callback) {
  return new EventListenerSignal().init(eventName, element, callback);
}

type EventListenerCallback = (this: Element, ev: Event) => any;

class EventListenerSignal<
  EventName extends keyof ElementEventMap,
  Target extends MaybeSignal<Element>,
  Callback extends MaybeSignal<EventListenerCallback>,
> extends ReadonlySignal<void> {
  public async init(event: EventName, element: Target, callback: Callback) {
    let currentListener: EventListenerCallback | null = null;

    const fn = async () => {
      const targetElement = await unwrapSignal(element);
      const handler = await unwrapSignal(callback);

      if (currentListener) {
        targetElement.removeEventListener(event, currentListener);
      }

      targetElement.addEventListener(event, callback as any);

      currentListener = callback as any;
    };

    fn();

    if (isSignal(element)) {
      element.subscribe(fn);
    }

    if (isSignal(callback)) {
      callback.subscribe(fn);
    }
  }
}
