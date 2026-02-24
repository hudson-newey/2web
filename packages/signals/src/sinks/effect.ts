import { execCallback } from "../utils/execCallback.ts";

type UseEffectReducer = () => unknown;

export function effect(callback: UseEffectReducer): void {
  const { dependencies } = execCallback(callback);

  for (const subscriber of dependencies) {
    subscriber.subscribe(callback);
  }
}
