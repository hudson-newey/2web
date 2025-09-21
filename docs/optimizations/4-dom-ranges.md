# DOM Ranges

The 2Web compiler automatically creates DOM Ranges of frequently mutated
subtrees.

This allows partial mutation updates to isolated parts of the DOM tree without
re-evaluating the entire document.
This is similar to React's "diffing" method, but better optimized because it
uses browser-native API's instead of a vDOM.

Note that unlike React, 2Web updates are targeted.
Therefore DOM ranges are only used when **dynamically** creating or removing
elements.

[DOM Ranges (MDN)](https://developer.mozilla.org/en-US/docs/Web/API/Range)
