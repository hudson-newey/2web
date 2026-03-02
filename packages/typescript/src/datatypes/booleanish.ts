import type { Variable } from "./memory";

export type Truthy<T> = T extends object ? T : never;
export type Falsy<T> = T extends object ? never : T;

export type FalsyValues = false | 0 | 0n | "" | null | undefined;
export type TruthyValues = Exclude<Variable, FalsyValues>;
