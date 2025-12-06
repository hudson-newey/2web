import { ReadonlySignal } from "../readonlySignal";

interface ResourceSignalOptions extends RequestInit {
  /**
   * @default 5000
   */
  timeout?: number;

  /**
   * @default false
   */
  streaming?: boolean;

  abortSignal?: AbortSignal;
}

type Resource<T> = Response & T;

class InternalResourceSignal<
  T,
  ResourceUrl extends string
> extends ReadonlySignal<Resource<T> | null> {
  public constructor(
    private readonly url: ResourceUrl,
    private readonly options?: ResourceSignalOptions
  ) {
    super(null);
  }

  public async init() {
    await this.refreshResource();
  }

  private async refreshResource() {
    const response = await fetch(this.url, this.options);
    this.value = response as any;
  }
}

/**
 * @description
 * A signal that fetches a resource and updates whenever the parameters change.
 *
 * If the resource is a WebSocket, it will automatically connect and emit a
 * new value whenever a message is received.
 * If one of the dependencies changes, a new connection will be established.
 */
export const ResourceSignal = new Proxy(InternalResourceSignal, {
  async construct(target, args) {
    const data = await Reflect.apply(target, this, args).init();
    return new target(data);
  },
});
