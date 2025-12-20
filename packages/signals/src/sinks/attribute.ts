import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export function attribute<
  Target extends HTMLElement,
  Attribute extends string,
  T
>(
  element: MaybeSignal<Target>,
  name: MaybeSignal<Attribute>,
  signal: Signal<T>
) {
  const target = unwrapSignal(element);
  const attributeName = unwrapSignal(name);

  signal.subscribe((value) => {
    if (value === null || value === undefined) {
      updateDom(() => {
        target.removeAttribute(attributeName);
      });

      return;
    }

    updateDom(() => {
      target.setAttribute(attributeName, String(value));
    });
  });
}
