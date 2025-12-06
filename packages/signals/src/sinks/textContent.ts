import type { Signal } from "../signal";
import { isSignal } from "../utils/isSignal";

/**
 * @description
 * A signal sink that safely sets the `textContent` of a DOM node to the value
 * of a signal.
 * This signal automatically escapes problematic characters to prevent XSS.
 */
export function textContent<const T>(
  node: Node | Signal<Node>,
  signal: Signal<T>
) {
  // TODO: Fix this horrible typing
  const target: Node = isSignal(node) ? (node.value as Node) : (node as Node);

  signal.subscribe((value) => {
    target.textContent = String(value);
  });
}
