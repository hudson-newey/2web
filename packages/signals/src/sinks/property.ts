import type { Signal } from "../signal";

export function property<
  Target extends HTMLElement,
  Property extends keyof Target,
  T extends Target[Property],
>(element: Target, name: Property, signal: Signal<T>) {
  signal.subscribe((value) => {
    element[name] = value;
  });
}
