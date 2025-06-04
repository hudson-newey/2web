# Code splitting

You may have noticed that sometimes the 2Web compiler does not always output a
single file.
This is an implementation detail that you should not have to worry about as I
believe that different rendering strategies should be deiced by the compiler
instead of the programmer, and that all code should be the same regardless of
the rendering strategy.

However, understanding code splitting in 2Web is an important optimization that
is automatically performed.

During code-splitting, the compiler has the following goals:

- The browser should only have to parse one file on initial page load.
  - This means that content on the critical rendering path should all be
    contained within the single `.html` file served to the user.
- Reactivity hydration should **not** change the information on the page.
  - When we hydrate the page with reactivity. E.g. "when I click this button xyz
    should occur" nothing about the page content should change.
- Loading reactivity should be non-blocking, and the page should be navigable
  before reactivity is loaded. Similar to [qwik](https://qwik.dev/), reactivity
  should be lazy loaded. This is beneficial because typical users will not
  immediately start interacting with content once it is visible. They will
  typically have to visually find the button they want to click, and move their
  meaty hands to click the mouse button. By delaying reactivity hydration, the
  user can visually scan the page to find the button they want to click, and by
  the time that they have found the button/element they want, reactivity should
  have loaded.
- Styling should be visible on first load so that there is no flash of unstyled
  content.
