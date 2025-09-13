import type { Signal } from "../signal";

export function pipe<T>(signal: Signal<T>, ...fns: Array<(value: T) => T>) {}
