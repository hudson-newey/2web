export class Timeline {
  private intervalRef: ReturnType<typeof setInterval> | null = null;

  constructor(private readonly callback: () => void) {}

  public get isRunning(): boolean {
    return this.intervalRef !== null;
  }

  public start() {
    this.intervalRef = setInterval(this.callback, 1);
  }

  public stop() {
    if (!this.intervalRef) {
      console.warn("Attempted to stop a timeline that is not running.");
      return;
    }

    clearInterval(this.intervalRef);
    this.intervalRef = null;
  }
}
