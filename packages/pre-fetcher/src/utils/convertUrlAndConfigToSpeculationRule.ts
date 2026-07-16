import type { PrefetchConfig } from "../config";
import type { SpeculationRule } from "../speculationScriptElement/speculationRule";

export function convertUrlAndConfigToSpeculationRule(
  url: string,
  config: Readonly<PrefetchConfig>,
): SpeculationRule {
  return {
    ...config,
    urls: [url],
  };
}
