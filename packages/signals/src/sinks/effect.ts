import { Signal } from "../signal.ts";

type UseEffectReducer = () => unknown;

export function effect<T>(callback: UseEffectReducer, subscribers: Signal<T>[]) {
  for (const subscriber of subscribers) {
    subscriber.subscribe(callback);
  }
}
