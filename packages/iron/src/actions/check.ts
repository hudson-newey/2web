/**
 * @description
 * Performs an action twice and asserts that the results are the same.
 */
export function check(fn: () => unknown): void {
  const result1 = fn();
  const result2 = fn();

  if (result1 !== result2) {
    throw new Error(`Check failed: '${result1}' !== '${result2}'`);
  }
}
