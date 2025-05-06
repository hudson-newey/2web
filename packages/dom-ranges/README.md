# 2Web Kit - Dom Ranges

A package to automatically create
[dom ranges](https://developer.mozilla.org/en-US/docs/Web/API/Range) from
compiled 2web code.

---

The execution of the dom-ranges script will be deferred until everything on the
page has finished rendering / execution and the DOM is in a stable state.

Once the DOM is in a stable state, this package will create dom-range fragments
for more efficient DOM updates.
