import { unwrapSignal } from "../..";
import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import type { MaybeSignal } from "../utils/unwrapSignal";

interface WithHref extends Node {
  href: string;
}

export function href<const T>(node: MaybeSignal<WithHref>, signal: Signal<T>) {
  const target = unwrapSignal(node);

  signal.subscribe((value) => {
    updateDom(() => {
      target.href = String(value);
    });
  });
}
