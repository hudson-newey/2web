# 2Web Kit - Pre-fetcher

Improves the performance of your web application by pre-fetching/pre-rendering
links / locations.

## Usage

```html
<script type="module">
  import { bootstrapLinkPrefetch, prefetch } from "@two-web/kit/pre-fetcher";

  // Call boostrapLinkPrefetch() to automatically attach event listeners to all
  // in-site links.
  bootstrapLinkPrefetch();

  // You can pass a config object into the bootstrap call
  bootstrapLinkPrefetch({
    where: {
      and: [
        { selector_matches: ".product-link" },
        { href_matches: "/user/*" },
        { not: { href_matches: "/logout" } } },
      ],
    },
  });

  // If you want to manually pre-fetch a url, you can directly call prefetch()
  // Calling prefetch() directly does NOT require you to bootstrap first.
  prefetch("/users/*");
  prefetch("https://www.google.com");
</script>
```
