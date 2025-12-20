import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal sink that safely sets the `textContent` of a DOM node to the value
 * of a signal.
 * This signal automatically escapes problematic characters to prevent XSS.
 */
export function textContent<const T>(
  node: MaybeSignal<Node>,
  signal: Signal<T>
) {
  const target = unwrapSignal(node);

  signal.subscribe((value) => {
    updateDom(() => {
      target.textContent = String(value);
    });
  });
}
