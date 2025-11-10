/**
 * @description
 * Requires you to handle error cases.
 */
export function tryCatch<T>(fn: () => T): T | Error {
  try {
    return fn();
  } catch (error) {
    if (error instanceof Error) {
      return error;
    }

    return new Error(String(error));
  }
}
