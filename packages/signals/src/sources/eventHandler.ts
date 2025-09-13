import { ReadonlySignal } from "../readonlySignal";

export type EventHandlerReducer<EventType extends Event, T> = (
  event: EventType,
  currentValue: T | null,
) => T | null;

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
export class EventHandler<
  T,
  EventType extends Event,
> extends ReadonlySignal<T | null> {
  private readonly reducer: EventHandlerReducer<EventType, T>;

  public constructor(reducer: EventHandlerReducer<EventType, T>) {
    super(null);
    this.reducer = reducer;
  }

  // handleEvent is a special method that allows this class to be passed
  // directly to addEventListener as the event listener.
  // https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#the_event_listener_callback
  public handleEvent(event: EventType) {
    this.value = this.reducer(event, this.value);
  }
}
