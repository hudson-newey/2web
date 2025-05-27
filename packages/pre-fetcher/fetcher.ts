const prefetchedLinks = new Set();

export function prefetchLink(target: HTMLAnchorElement): void {
  // Pre-fetching remote origins is a security and privacy risk.
  // Therefore, we only pre-fetch targets that point to the same origin.
  if (!isSameOrigin(target)) {
    return;
  }

  const hrefTarget = target.href;
  if (prefetchedLinks.has(hrefTarget)) {
    return;
  }

  // We add the href target to the completed heap before requesting the link
  // so that we don't get a race condition where multiple fetch requests are
  // made in async before we append to the completed heap.
  prefetchedLinks.add(hrefTarget);

  // Notice that this fetch is not awaited so that it is performed in the
  // microtask queue.
  fetch(hrefTarget);
}

function isSameOrigin(target: HTMLAnchorElement): boolean {
  const href = target.href;
  const linkUrl = new URL(href, window.location.origin);
  const currentUrl = window.location;

  return (
    linkUrl.protocol === currentUrl.protocol &&
    linkUrl.host === currentUrl.host &&
    linkUrl.port === currentUrl.port
  );
}
