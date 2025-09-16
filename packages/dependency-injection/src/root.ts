import { injectAdvertisementEvent, injectRequestEvent } from "./events";
import { isInjectableRequest } from "./inject";
import type { Injectable, InjectionConstructor } from "./injectable";
import { isProviderAdvertisement } from "./provider";

export class InjectionRoot {
  private readonly providers = new Map<InjectionConstructor, Injectable>();
  private readonly rootElement: HTMLElement;

  public constructor(rootElement: HTMLElement = document.body) {
    this.rootElement = rootElement;

    this.attachModule();
  }

  private attachModule(): void {
    this.rootElement.addEventListener(injectAdvertisementEvent, (event) => {
      if (!isProviderAdvertisement(event)) {
        console.warn(`Received invalid provider advertisement. Expected "event.detail" to be defined.`);
        return;
      }

      this.providers.set(event.detail.ctor, event.detail.provider);

      event.stopPropagation();
    });

    this.rootElement.addEventListener(injectRequestEvent, (event) => {
      if (!isInjectableRequest(event)) {
        return;
      }

      const provider = this.providers.get(event.detail.ctor);
      if (provider) {
        event.detail.resolver(provider);
      } else {
        console.warn(`No provider found for requested injectable. "${event.detail.ctor.name}"`);
      }

      event.stopPropagation();
    });
  }
}
