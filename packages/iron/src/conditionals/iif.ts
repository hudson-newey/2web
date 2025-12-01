export interface IifOptions<T, Then = T, Else = T, Catch = T> {
  then?: () => Then;
  else?: () => Else;
  catch?: (error: unknown) => Catch;
}

/**
 * @description
 * When creating high-reliability software, type guards can cause the program to
 * hard fail if an unexpected type is reached.
 *
 * This conditional handles errors gracefully without causing the program to
 * crash.
 *
 * @returns
 * A boolean indicating whether the condition was met.
 * If `then` or `else` callbacks are provided ,their return values will be used
 * instead of the default `true` or `false`.
 * Note that if the `then` or `else` callbacks return `undefined`, the default
 * boolean values will be used.
 */
export function iif<T = void>(
  condition: () => Promise<T> | T
): (options: IifOptions<T>) => Promise<T | boolean> {
  return async (options: IifOptions<T>) => {
    try {
      const result = await condition();
      if (result) {
        return options.then?.() ?? true;
      } else {
        return options.else?.() ?? false;
      }
    } catch (error) {
      if (options.catch) {
        return options.catch(error);
      }

      // If no catch handler is provided, we throw the error to the caller.
      throw error;
    }
  };
}
