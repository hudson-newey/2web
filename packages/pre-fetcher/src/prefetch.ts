import type { PrefetchConfig } from "./config";
import { SpeculationScript } from "./speculationScriptElement/speculationScriptElement";
import { mergeDefaultConfig } from "./utils/mergeDefaultConfig";

export function prefetch(
  target: string | HTMLAnchorElement,
  config: PrefetchConfig = {},
) {
  const href: string =
    target instanceof HTMLAnchorElement ? target.href : target;

  const mergedConfig = mergeDefaultConfig(config);
  SpeculationScript.addLink(href, mergedConfig);
}
