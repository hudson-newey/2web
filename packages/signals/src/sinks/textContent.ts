import type { Signal } from "../signal";

/**
 * @description
 * A signal sink that safely sets the `textContent` of a DOM node to the value
 * of a signal.
 * This signal automatically escapes problematic characters to prevent XSS.
 */
export function textContent<T>(node: Node, signal: Signal<T>) {
  signal.subscribe((value) => {
    node.textContent = String(value);
  });
}
