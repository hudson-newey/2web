import { ReadonlySignal } from "../readonlySignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal that fetches a resource and updates whenever the parameters change.
 *
 * If the resource is a WebSocket, it will automatically connect and emit a
 * new value whenever a message is received.
 * If one of the dependencies changes, a new connection will be established.
 */
export function resource<T, ResourceUrl extends string>(
  urlInput: MaybeSignal<ResourceUrl>,
  optionsInput?: MaybeSignal<ResourceSignalOptions>
) {
  return new ResourceSignal<T, ResourceUrl>(urlInput, optionsInput);
}

export interface ResourceSignalOptions extends RequestInit {
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
  private readonly url: ResourceUrl;
  private readonly options?: ResourceSignalOptions;

  public constructor(
    urlInput: MaybeSignal<ResourceUrl>,
    optionsInput?: MaybeSignal<ResourceSignalOptions>
  ) {
    super(null);

    this.url = unwrapSignal(urlInput);

    if (optionsInput !== undefined) {
      this.options = unwrapSignal(optionsInput);
    }
  }

  public async init() {
    await this.refreshResource();
  }

  private async refreshResource() {
    const response = await fetch(this.url, this.options);
    this.value = response as any;
  }
}

const ResourceSignal = new Proxy(InternalResourceSignal, {
  async construct(target, args) {
    const data = await Reflect.apply(target, this, args).init();
    return new target(data);
  },
});
