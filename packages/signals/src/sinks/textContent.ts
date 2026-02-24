import { updateDom } from "../../../_shared/updateDom";
import type { Signal } from "../signal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal sink that safely sets the `textContent` of a DOM node to the value
 * of a signal.
 * This signal automatically escapes problematic characters to prevent XSS.
 */
export async function textContent<const T>(
  node: MaybeSignal<Node>,
  signal: Signal<T>
) {
  const fn = (value: T | null) => {
    updateDom(async () => {
      const target = await unwrapSignal(node);
      target.textContent = String(value);
    });
  };

  signal.subscribe(fn);

  if (isSignal(node)) {
    node.subscribe(() => fn(signal.value));
  }
}
