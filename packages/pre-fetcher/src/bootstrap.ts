import { defaultPrefetchConfig, type PrefetchConfig } from "./config";
import { Prefetcher } from "./fetcher";

export function bootstrapLinkPrefetch(config: PrefetchConfig) {
  const mergedConfig = Object.assign({}, defaultPrefetchConfig, config);
  const prefetcher = new Prefetcher(mergedConfig);

  const anchorElements = document.getElementsByTagName("a");
  for (const anchorTarget of anchorElements) {
    prefetcher.observe(anchorTarget);
  }

  const documentMutation = new MutationObserver((mutationList) => {
    for (const mutation of mutationList) {
      for (const node of mutation.addedNodes) {
        if (node instanceof HTMLAnchorElement) {
          prefetcher.observe(node);
        }
      }

      // We have to correctly remove event listeners so that there are no
      // hanging references to the observed elements which might cause a memory
      // leak.
      for (const node of mutation.removedNodes) {
        if (node instanceof HTMLAnchorElement) {
          prefetcher.disconnect(node);
        }
      }
    }
  });

  documentMutation.observe(document.body, {
    subtree: true,
    childList: true,
    attributes: true,
    attributeFilter: ["href", "target", "rel"],
  });
}
