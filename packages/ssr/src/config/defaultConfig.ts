import type { SsrConfig } from "./config";

export const defaultConfig = Object.freeze({
    // Should match the default port for the SPA dev server to reduce confusion.
    port: 2000,
}) satisfies Required<SsrConfig>;

export function mergeDefaultConfig(config: Readonly<Partial<SsrConfig>>) {
  return Object.assign({}, config, defaultConfig);
}
