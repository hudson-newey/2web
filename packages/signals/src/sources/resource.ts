import { ReadonlySignal } from "../readonlySignal";

interface ResourceSignalOptions<T> {
  /**
   * @default 5000
   */
  timeout?: number;

  /**
   * @default false
   */
  streaming?: boolean;

  fetchOptions?: RequestInit;

  abortSignal?: AbortSignal;
}

/**
 * @description
 * A signal that fetches a resource and updates whenever the parameters change.
 *
 * If the resource is a WebSocket, it will automatically connect and emit a
 * new value whenever a message is received.
 * If one of the dependencies changes, a new connection will be established.
 */
export class ResourceSignal<T> extends ReadonlySignal<T> {}
