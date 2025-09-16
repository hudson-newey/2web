import { injectAdvertisementEvent } from "./events";
import type { Injectable, InjectionConstructor } from "./injectable";

export interface ProviderAdvertisement {
  ctor: InjectionConstructor;
  provider: Injectable;
};

export function isProviderAdvertisement(
  event: unknown,
): event is CustomEvent<ProviderAdvertisement> {
  return true;
}

export function provide(
  token: InjectionConstructor,
  provider: Injectable,
  from: HTMLElement = document.body
): void {
  const event = new CustomEvent<ProviderAdvertisement>(injectAdvertisementEvent, {
    detail: { ctor: token, provider },
  });

  from.dispatchEvent(event);
}
