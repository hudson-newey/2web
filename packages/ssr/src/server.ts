import { createSsrHandler } from "./renderer/handler.ts";
import type { SsrConfig } from "./config/config.ts";
import { mergeDefaultConfig } from "./config/defaultConfig.ts";

export function runServer(userConfig: Readonly<Partial<SsrConfig>>) {
  const config = mergeDefaultConfig(userConfig);
  const handler = createSsrHandler(config);

  Deno.serve(handler);
}
