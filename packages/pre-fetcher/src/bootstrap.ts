import type { PrefetchConfig } from "./config";
import { SpeculationScript } from "./speculationScriptElement/speculationScriptElement";

export function bootstrapLinkPrefetch(config: Readonly<PrefetchConfig>) {
  SpeculationScript.addRule({
    eagerness: "eager",
    where: {
      selector_matches: "a",
    },
    ...config,
  });
}
