/**
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#speculation_rules_json_representation}
 */
export type SpeculationRule = UrlSpeculationRule | WhereSpeculationRule;

/**
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#urls}
 */
interface UrlSpeculationRule extends BaseSpeculationRule {
  urls: string[];
}

interface WhereSpeculationRule {
  where: WhereSpeculationRuleCondition;
}

/**
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#where}
 */
type WhereSpeculationRuleCondition = BaseSpeculationRule &
  (
    | { href_matches: string }
    | { relative_to: string }
    | { selector_matches: string }
    | { and: RecursiveWhereSpeculationRule }
    | { not: RecursiveWhereSpeculationRule }
    | { or: RecursiveWhereSpeculationRule }
  );

// Inside of "and", "not", and "or" blocks, you cannot use the "relative_to"
// rule.
type RecursiveWhereSpeculationRule = Exclude<
  WhereSpeculationRule,
  "relative_to"
>;

/**
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/speculationrules#speculation_rules_json_representation}
 */
export interface BaseSpeculationRule {
  eagerness?: "immediate" | "eager" | "moderate" | "conservative";
  expects_no_vary_search?: boolean;
  referrer_policy?: ReferrerPolicy;
  relative_to?: "document" | "ruleset";
  requires?: "anonymous-client-ip-when-cross-origin";
  tag?: string;
  target_hint?: "_blank" | "_self";
}

type ReferrerPolicy =
  | "strict-origin-when-cross-origin"
  | "same-origin"
  | "strict-origin"
  | "no-referrer";
