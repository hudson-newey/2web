export type Xor<T1, T2> = (T1 & Exclude<T2, T1>) | (Exclude<T1, T2> & T2);
