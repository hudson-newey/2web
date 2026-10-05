import type { PrefetchConfig } from "../config.ts";

export const defaultPrefetchConfig = Object.freeze({
  eagerness: "eager",
}) satisfies PrefetchConfig;

export function mergeDefaultConfig(config: Readonly<PrefetchConfig>) {
  return Object.assign({}, config, defaultPrefetchConfig);
}
