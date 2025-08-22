import type { PrefetchConfig } from "./config";

type ElementEvent = keyof HTMLElementEventMap;

export class Prefetcher {
  // We use a WeakMap so that we know what links have event listeners attached
  // without creating a hard reference so that we don't create a memory leak where
  // the element nodes can't be cleaned up.
  private readonly listenedElements = new WeakSet();

  // We use a set of prefetched link urls instead of HTMLAnchorElement
  // references so that if the link href changes, the prefetch state is
  // invalidated.
  private readonly prefetchedLinks = new Set();

  // We use a bound event handler so that when we go to remove it, we are 100%
  // certain that it has the same heap reference.
  private readonly prefetchHandler = this.prefetchLink.bind(this);

  public constructor(private readonly config: PrefetchConfig) {}

  private get observedEvent(): ElementEvent {
    const triggerMap = new Map<PrefetchConfig["trigger"], ElementEvent>([
      ["click", "pointerdown"],
      ["hover", "pointerenter"],
      ["focus", "focus"],
    ]);

    // Because the config is user-facing, it can be used in untyped environments
    // where they might input an invalid trigger like a boolean.
    // In this case, "undefined" will be returned from the trigger map, and we
    // will default to "pointerenter"
    const mappedEvent = triggerMap.get(this.config.trigger);
    if (mappedEvent === undefined) {
      console.error(
        `Unknown prefetch trigger: '${this.config.trigger}'. Defaulting to 'hover'.`,
      );
      return triggerMap.get("hover") as ElementEvent;
    }

    return mappedEvent;
  }

  public prefetchLink(event: Event): void {
    const target = event.target as HTMLAnchorElement;

    // Pre-fetching remote origins is a security and privacy risk.
    // Therefore, we only pre-fetch targets that point to the same origin.
    if (!this.isSameOrigin(target)) {
      return;
    }

    const hrefTarget = target.href;
    if (this.prefetchedLinks.has(hrefTarget)) {
      return;
    }

    // We add the href target to the completed heap before requesting the link
    // so that we don't get a race condition where multiple fetch requests are
    // made in async before we append to the completed heap.
    this.prefetchedLinks.add(hrefTarget);

    // Notice that this fetch is not awaited so that it is performed in the
    // microtask queue.
    fetch(hrefTarget, this.config);
  }

  public observe(target: HTMLAnchorElement): void {
    if (this.listenedElements.has(target)) {
      return;
    }

    target.addEventListener(this.observedEvent, this.prefetchHandler);

    this.listenedElements.add(target);
  }

  public disconnect(target: HTMLAnchorElement): void {
    if (!this.listenedElements.has(target)) {
      return;
    }

    target.removeEventListener(this.observedEvent, this.prefetchHandler);

    this.listenedElements.delete(target);
  }

  private isSameOrigin(target: HTMLAnchorElement): boolean {
    const href = target.href;
    const linkUrl = new URL(href, window.location.origin);
    const currentUrl = window.location;

    return (
      linkUrl.protocol === currentUrl.protocol &&
      linkUrl.host === currentUrl.host &&
      linkUrl.port === currentUrl.port
    );
  }
}
