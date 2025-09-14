export type PrefetchTrigger = "hover" | "focus" | "click";

export interface PrefetchConfig extends RequestInit {
  /** @default "hover" */
  trigger?: PrefetchTrigger;
}

export const defaultPrefetchConfig = Object.freeze({
  trigger: "hover",
  priority: "low",
}) satisfies PrefetchConfig;
