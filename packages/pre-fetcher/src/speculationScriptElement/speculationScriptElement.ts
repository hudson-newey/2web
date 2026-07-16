import type { PrefetchConfig } from "../config";
import { convertUrlAndConfigToSpeculationRule } from "../utils/convertUrlAndConfigToSpeculationRule";
import type { SpeculationRule } from "./speculationRule";

/**
 * @summary
 * A programatic way to interact with a document speculation rule element
 * without needing to constantly compute script content + handle de-duping and
 * other edge cases.
 *
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules}
 */
export class SpeculationScript {
  public static addLink(url: string, config: Readonly<PrefetchConfig>) {
    SpeculationScript.addRule(
      convertUrlAndConfigToSpeculationRule(url, config),
    );
  }

  public static addRule(rule: Readonly<SpeculationRule>) {
    // Unfortunately, we need to re-create a script element every time we want
    // to invalidate prefetch rules.
    const element = document.createElement("script");
    element.type = "speculationrules";
    element.textContent = JSON.stringify({ prerender: [rule] });
    document.body.append(element);
  }
}
