import { isSignal, unwrapSignal } from "../..";
import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import type { MaybeSignal } from "../utils/unwrapSignal";

interface WithHref extends Node {
  href: string;
}

export async function href<const T>(node: MaybeSignal<WithHref>, signal: Signal<T>) {
  const fn = async (value: T | null) => {
    const target = await unwrapSignal(node);

    updateDom(() => {
      target.href = String(value);
    });
  }

  signal.subscribe(fn);

  if (isSignal(node)) {
    node.subscribe(() => fn(signal.value));
  }
}
