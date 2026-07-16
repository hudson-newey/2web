import type { PrefetchConfig } from "../config";

export const defaultPrefetchConfig = Object.freeze({
  eagerness: "eager",
}) satisfies PrefetchConfig;

export function mergeDefaultConfig(config: Readonly<PrefetchConfig>) {
  return Object.assign({}, config, defaultPrefetchConfig);
}
