console.warn(
  "You are importing the @2web/kit barrel file. This is not recommended and " +
    "can lead to larger bundle sizes, and reduce performance. If performance " +
    "a concern, please cherry pick the individual packages you need instead.",
);

export * from "./animations/index.ts";
export * from "./iron/index.ts";
export * from "./named-strings/index.ts";
export * from "./pre-fetcher/index.ts";
export * from "./route-guards/index.ts";
export * from "./ssr/index.ts";
export * from "./typescript/index.ts";
export * from "./2js/index.ts";
