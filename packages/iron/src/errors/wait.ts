export interface WaitOptions {
  catch?: (error: unknown) => unknown;
}

/**
 * @description
 * An awaitable function to unwrap promises while handling their rejections.
 *
 * This is needed because a normal async function does have a `catch` callback, but
 * this cannot be easily used with the `await` keyword.
 */
export async function wait<T>(
  fn: () => T | Promise<T>,
  options: WaitOptions = {}
): Promise<
  ReturnType<typeof fn> | ReturnType<NonNullable<WaitOptions["catch"]>>
> {
  try {
    return await fn();
  } catch (error) {
    if (options.catch) {
      return options.catch(error);
    }

    throw error;
  }
}
