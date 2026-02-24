import { ReadonlySignal } from "../readonlySignal";

/**
 * A signal that can also be passed directly to addEventListener.
 *
 * @example
 * ```ts
 * import { EventHandler } from "@two-web/kit/signals";
 *
 * const target = document.getElementById("counter");
 * const clickHandler = new EventHandler((event, value) => {
 *   const count = value + 1;
 *
 *   event.target.textContent = count;
 *   return count;
 * });
 *
 * target.addEventListener("click", clickHandler);
 * ```
 */
export function eventHandler<T, EventType extends Event>(
  initialValue: T | null,
  reducer: EventHandlerReducer<EventType, T>
): EventHandler<T, EventType> {
  return new EventHandler<T, EventType>(initialValue, reducer);
}

class EventHandler<
  T,
  EventType extends Event
> extends ReadonlySignal<T | null> {
  public constructor(
    initialValue: T | null,
    private readonly reducer: EventHandlerReducer<EventType, T>
  ) {
    super(initialValue);
  }

  // handleEvent is a special method that allows this class to be passed
  // directly to addEventListener as the event listener.
  // https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#the_event_listener_callback
  public handleEvent(event: EventType) {
    this._internalSet(this.reducer(event, this.value));
  }
}

type EventHandlerReducer<EventType extends Event, T> = (
  event: EventType,
  currentValue: T | null
) => T | null;
