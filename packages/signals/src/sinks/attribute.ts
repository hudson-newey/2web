import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export async function attribute<
  Target extends HTMLElement,
  Attribute extends string,
  T
>(
  element: MaybeSignal<Target>,
  name: MaybeSignal<Attribute>,
  signal: Signal<T>
) {
  const fn = async (value: T | null) => {
    const target = await unwrapSignal(element);
    const attributeName = await unwrapSignal(name);

    if (value === null || value === undefined) {
      updateDom(() => {
        target.removeAttribute(attributeName);
      });

      return;
    }

    updateDom(() => {
      target.setAttribute(attributeName, String(value));
    });
  };

  signal.subscribe(fn);

  if (isSignal(element)) {
    element.subscribe(() => fn(signal.value));
  }

  if (isSignal(name)) {
    name.subscribe(() => fn(signal.value));
  }
}
