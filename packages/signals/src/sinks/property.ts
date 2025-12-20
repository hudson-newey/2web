import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export function property<
  Target extends HTMLElement,
  Property extends keyof Target,
  T extends Target[Property]
>(
  element: MaybeSignal<Target>,
  name: MaybeSignal<Property>,
  signal: Signal<T>
) {
  const target = unwrapSignal(element);
  const propName = unwrapSignal(name);

  signal.subscribe((value) => {
    updateDom(() => {
      target[propName] = value;
    });
  });
}
