import type { InjectionConstructor } from "./injectable";

export interface InjectableRequest {
  ctor: InjectionConstructor;
  resolver: (injectable: any) => void;
}

export function isInjectableRequest(
  event: unknown
): event is CustomEvent<InjectableRequest> {
  return true;
}

export async function inject(
  token: InjectionConstructor,
  from: HTMLElement = document.body
): Promise<any> {
  return new Promise<any>((resolve) => {
    const event = new CustomEvent<InjectableRequest>("inject-request", {
      detail: { ctor: token, resolver: resolve },
      bubbles: true,
      composed: true,
    });

    from.dispatchEvent(event);
  });
}
