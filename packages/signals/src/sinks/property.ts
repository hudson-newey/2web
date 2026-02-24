import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

export async function property<
  Target extends HTMLElement,
  Property extends keyof Target,
  T extends Target[Property]
>(
  element: MaybeSignal<Target>,
  name: MaybeSignal<Property>,
  signal: Signal<T>
) {
  const fn = async (value: T | null) => {
    const target = await unwrapSignal(element);
    const propName = await unwrapSignal(name);

    updateDom(() => {
      (target[propName] as unknown) = value;
    });
  }

  signal.subscribe(fn);

  if (isSignal(element)) {
    element.subscribe(() => fn(signal.value));
  }

  if (isSignal(name)) {
    name.subscribe(() => fn(signal.value));
  }
}
