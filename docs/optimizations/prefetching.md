# 2Web Resource Pre-fetching

2Web improves rendering performance by pre-fetching resources such as css
background images.

While modern browsers have started pre-fetching resources such as `<img>`
sources, support is limited
(e.g. at the time of writing, css `background-image`'s do not work), and the
html parser needs to get to the `<img>` element node before pre-fetching.

For large resources, 2Web will automatically add `<link rel="prefetch">` tags
to your documents head if the external resource is in the users initial
viewport.

Note that if there are a large number of resources that would benefit from
pre-fetching, 2Web will limit prefetching to the largest _6_ elements that are
in the users initial viewport.

We will only prefetch resources in the initial viewport, because there is no
need to prefetch a large resource at the bottom of the page that can be lazy
loaded.
