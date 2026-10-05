import type { PrefetchConfig } from "../config.ts";
import type { SpeculationRule } from "../speculationScriptElement/speculationRule.ts";

export function convertUrlAndConfigToSpeculationRule(
  url: string,
  config: Readonly<PrefetchConfig>,
): SpeculationRule {
  return {
    ...config,
    urls: [url],
  };
}
