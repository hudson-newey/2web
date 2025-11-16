export interface IifOptions<T> {
  then?: () => T;
  else?: () => T;
  catch?: (error: unknown) => T;
}

/**
 * @description
 * When creating high-reliability software, type guards can cause the program to
 * hard fail if an unexpected type is reached.
 *
 * This conditional handles errors gracefully without causing the program to
 * crash.
 */
export function iif<T = void>(
  condition: () => boolean
): (options: IifOptions<T>) => T | undefined {
  return (options: IifOptions<T>) => {
    try {
      const result = condition();
      if (result) {
        return options.then?.();
      } else {
        return options.else?.();
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
