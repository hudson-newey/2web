export type Truthy<T> = T extends object ? T : never;
export type Falsy<T> = T extends object ? never : T;
