# 2Web Chrome Extension

A Chrome DevTools extension for inspecting 2web development builds.

The extension adds a "2web" panel to the Chrome DevTools. The panel reads the
`__2web.debug.json` file that the 2web compiler writes into the build output of
development builds (the file is omitted from production builds) and renders
three tabs:

- **Assets** — every file in the compiled build output (pages, stylesheets,
  scripts, and passthrough assets) with its type and size.
- **Reactivity** — the reactive variables (with their reactivity class),
  reactive properties (with the variables their reducer depends on), and
  reactive events (with their assignment sink and dependencies) in a table
  format.
- **Graph** — the relationship graph between variables, properties, and
  events. Clicking a node highlights everything it's connected to.

## Usage

1. Build the extension:

   ```sh
   pnpm build
   ```

2. Load the `dist/` directory as an unpacked extension through
   `chrome://extensions` (developer mode).

3. Open the DevTools of a page served from a 2web development build and switch
   to the "2web" panel. Use the "Refresh" button to re-read the debug file
   after recompiling.

## Architecture

- `public/devtools.html` / `public/devtools.js` — the devtools page that
  registers the "2web" panel.
- `index.html` + `src/script.ts` — the panel itself. The debug file is fetched
  from the inspected page's origin with `chrome.devtools.inspectedWindow.eval`,
  which doesn't require any manifest permissions.
