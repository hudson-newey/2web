import { createSsrHandler } from "./renderer/handler";
import express from "express";
import type { SsrConfig } from "./config/config";
import { mergeDefaultConfig } from "./config/defaultConfig";

export async function runServer(userConfig: Readonly<Partial<SsrConfig>>) {
  const config = mergeDefaultConfig(userConfig);
  const app = express();

  const handler = createSsrHandler(config);
  app.use(/(.*)/, handler);

  app.listen(config.port);
}
