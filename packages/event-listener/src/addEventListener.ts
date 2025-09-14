export interface EventListenerOptions extends AddEventListenerOptions {}

export function eventListener(
  target: EventTarget,
  type: string,
  listener: EventListenerOrEventListenerObject,
  options?: EventListenerOptions,
) {
  target.addEventListener(type, listener, options);
}
