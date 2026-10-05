import type { PrefetchConfig } from "./config.ts";
import { SpeculationScript } from "./speculationScriptElement/speculationScriptElement.ts";

export function bootstrapLinkPrefetch(config: Readonly<PrefetchConfig>) {
  SpeculationScript.addRule({
    eagerness: "eager",
    where: {
      selector_matches: "a",
    },
    ...config,
  });
}
