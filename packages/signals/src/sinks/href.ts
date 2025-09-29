import type { Signal } from "../signal";

interface WithHref extends Node {
  href: string;
}

export function href<T>(node: WithHref, signal: Signal<T>) {
  signal.subscribe((value) => {
    node.href = String(value);
  });
}
