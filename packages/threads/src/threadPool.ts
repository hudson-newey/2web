export function threadPool() {}

export function createThreadPool(size: number): ThreadPool {
  return ThreadPool.setSize(size);
}

class ThreadPool {
  // Using up all of the threads does not leave any room for the main thread,
  // the operating system, or other applications to breathe.
  // Therefore, by default, we only use half of the available threads.
  public static defaultSize = (navigator.hardwareConcurrency / 2) || 4;

  private readonly threads: Worker[];

  private constructor(size: number) {
    this.threads = Array.from({ length: size }, () => new Worker());
  }

  public static setSize(size: number): ThreadPool {
    return new ThreadPool(size);
  }

  public get workerCount(): number {
    return this.threads.length;
  }
}
