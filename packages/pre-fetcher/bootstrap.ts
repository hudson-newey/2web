import { prefetchLink } from "./fetcher";

const listenedElements = new WeakSet();

export function bootstrapLinkPrefetch() {
  const anchorElements = document.getElementsByTagName("a");
  for (const anchorTarget of anchorElements) {
    attachListener(anchorTarget);
  }

  const documentMutation = new MutationObserver((mutationList) => {
    for (const mutation of mutationList) {
      for (const node of mutation.addedNodes) {
        if (node instanceof HTMLAnchorElement) {
          attachListener(node);
        }
      }
    }
  });

  documentMutation.observe(document.body);
}

function attachListener(target: HTMLAnchorElement): void {
  if (listenedElements.has(target)) {
    return;
  }

  target.addEventListener("pointerenter", () => {
    prefetchLink(target);
  });

  listenedElements.add(target);
}
