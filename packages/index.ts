console.warn(
  "You are importing the 2Web kit barrel file. This is not recommended as it " +
    "can lead to larger bundle sizes and reduce performance. If performance " +
    "a concern, please cherry pick the individual packages you need instead.",
);

export * from "./animations/index.ts";
export * from "./browser-state/index.ts";
export * from "./iron/index.ts";
export * from "./dependency-injection/index.ts";
export * from "./event-listener/index.ts";
export * from "./hydration/index.ts";
export * from "./keyboard/index.ts";
export * from "./pre-fetcher/index.ts";
export * from "./route-guards/index.ts";
export * from "./signals/index.ts";
export * from "./ssr/index.ts";
export * from "./threads/index.ts";
export * from "./vdom/index.ts";

// We do not export the vite plugin here because it is not intended to be used
// within browser environments.
