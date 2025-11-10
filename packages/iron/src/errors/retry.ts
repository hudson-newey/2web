export interface RetryOptions {
  interval?: number;
  times?: number;
}

export async function retry<ReturnType>(
  fn: () => ReturnType | Promise<ReturnType>,
  options: RetryOptions,
) {
  const { interval = 1000, times = 3 } = options;

  let attempt = 0;
  while (attempt < times) {
    try {
      return await fn();
    }
    catch (error) {
      attempt++;
      if (attempt >= times) {
        throw error;
      }
      await new Promise((resolve) => setTimeout(resolve, interval));
    }
  }
}
