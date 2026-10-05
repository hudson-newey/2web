import type { PrefetchConfig } from "./config.ts";
import { SpeculationScript } from "./speculationScriptElement/speculationScriptElement.ts";
import { mergeDefaultConfig } from "./utils/mergeDefaultConfig.ts";

export function prefetch(
  target: string | HTMLAnchorElement,
  config: PrefetchConfig = {},
) {
  const href: string =
    target instanceof HTMLAnchorElement ? target.href : target;

  const mergedConfig = mergeDefaultConfig(config);
  SpeculationScript.addLink(href, mergedConfig);
}
