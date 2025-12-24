import { execCallback } from "../utils/execCallback.ts";

type UseEffectReducer = () => unknown;

export function effect(callback: UseEffectReducer): void {
  const { subscribers } = execCallback(callback);

  for (const subscriber of subscribers) {
    subscriber.subscribe(callback);
  }
}
