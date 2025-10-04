/**
 * A tagged template function that returns the input as-is, while preserving
 * structural typing.
 */
export function identityStringTemplate<const T extends string>(
  strings: TemplateStringsArray,
  ...values: any[]
): T {
  let result = "";

  for (let i = 0; i < strings.length; i++) {
    result += strings[i];
    if (i < values.length) {
      result += values[i];
    }
  }

  return result as T;
}
