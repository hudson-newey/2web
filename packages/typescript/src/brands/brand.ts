const brand = Symbol("__brand");

export type Brand<T = unknown, BrandName = string> = T & { [brand]: BrandName };

export function withBrand<const T, const BrandName extends string>(
  value: T,
  _brandName?: BrandName,
): Brand<T, BrandName> {
  return value as Brand<T, BrandName>;
}
