import { ReadonlySignal } from "../readonlySignal";

interface TimerSignalOptions {
  /**
   * @description
   * How often the timer should emit a new value.
   *
   * @units milliseconds
   * @default 16 // ms
   */
  interval?: number;

  /**
   * @description
   * The amount of time the timer should run before completing.
   *
   * @units milliseconds
   * @default Infinity
   */
  duration?: number;

  /**
   * @description
   * Whether the timer should start automatically upon creation.
   *
   * @default true
   */
  autoStart?: boolean;
}

export class TimerSignal extends ReadonlySignal<number> {
  private readonly options: Required<TimerSignalOptions>;
  private intervalId: ReturnType<typeof setInterval> | null = null;
  private elapsed = 0;

  public constructor({
    interval = 16,
    duration = Infinity,
    autoStart = true,
  }: TimerSignalOptions = {}) {
    super(0);

    this.options = { interval, duration, autoStart };

    if (this.options.autoStart) {
      this.start();
    }
  }

  public start(): this {
    const interval = this.options.interval;
    const duration = this.options.duration;

    this.intervalId = setInterval(() => {
      this.elapsed += interval;
      this.value = this.elapsed;

      if (this.elapsed >= duration) {
        this.stop();
      }
    }, interval);

    return this;
  }

  public stop(): this {
    if (this.intervalId !== null) {
      clearInterval(this.intervalId);
      this.intervalId = null;
    }

    return this;
  }
}
