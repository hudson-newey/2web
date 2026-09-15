// The dev server client runtime. The server injects this into every page it
// serves and notifies it over the /__2web_updates websocket:
//
//   - A rebuild that changed css assets hot swaps the stylesheets in place
//     (the page's state is kept). The fresh page is fetched to resolve the
//     (hashed) asset names that the rebuilt page references.
//   - Every other change (javascript, html markup, ...) reloads the page:
//     without bundler cooperation, re-executing javascript can't update the
//     module references the page already holds.
//
// The legacy plain "reload" message (sent by the action socket) reloads the
// page.

const socket = new WebSocket(`ws://${window.location.host}/__2web_updates`);

function reloadPage() {
  console.debug("[2web] Reloading page...");
  window.location.reload();
}

// Swaps every stylesheet link of the page for the version that the rebuilt
// page references (matched by their order of appearance, because the
// compiler names assets by a content hash that changes on every edit).
//
// Each swap happens once the new stylesheet has loaded, so the page never
// renders unstyled.
async function hotSwapStylesheets() {
  let freshPage;

  try {
    freshPage = await fetch(window.location.href, { cache: "no-store" });
  } catch (error) {
    console.error("[2web] Failed to fetch the updated page:", error);
    reloadPage();
    return;
  }

  const freshDocument = new DOMParser().parseFromString(
    await freshPage.text(),
    "text/html",
  );

  const freshLinks = [
    ...freshDocument.querySelectorAll('link[rel="stylesheet"]'),
  ];
  const currentLinks = [
    ...document.querySelectorAll('link[rel="stylesheet"]'),
  ];

  if (freshLinks.length !== currentLinks.length) {
    // The page's stylesheet structure changed; a reload is the only safe
    // update.
    reloadPage();
    return;
  }

  currentLinks.forEach((currentLink, index) => {
    const freshHref = freshLinks[index].getAttribute("href");

    if (!freshHref || currentLink.getAttribute("href") === freshHref) {
      // The stylesheet didn't change.
      return;
    }

    const next = currentLink.cloneNode();
    next.href = new URL(freshHref, window.location.href) + "?hmr=" + Date.now();
    next.addEventListener(
      "load",
      () => currentLink.remove(),
      { once: true },
    );
    next.addEventListener(
      "error",
      () => {
        console.error("[2web] Failed to load the updated stylesheet:", freshHref);
        currentLink.remove();
      },
      { once: true },
    );

    currentLink.parentNode.insertBefore(next, currentLink.nextSibling);
    console.debug("[2web] Hot swapped stylesheet:", freshHref);
  });
}

async function handleAssets(assets) {
  // A css change hot swaps the stylesheets. The page reload is skipped when
  // the css hot swap covered every changed asset, so that editing a
  // stylesheet keeps the page's state.
  const cssChanged = assets.some((asset) => asset.endsWith(".css"));
  const reload = assets.some((asset) => !asset.endsWith(".css"));

  if (cssChanged) {
    await hotSwapStylesheets();
  }

  if (reload) {
    reloadPage();
  }
}

socket.onmessage = (event) => {
  if (event.data === "reload") {
    reloadPage();
    return;
  }

  try {
    const message = JSON.parse(event.data);

    if (message.type === "assets" && Array.isArray(message.assets)) {
      handleAssets(message.assets);
    }
  } catch {
    // Not a message this client understands.
  }
};

socket.onclose = () => {
  console.error("[2web] Dev server connection lost.");
};

socket.onerror = (err) => {
  throw new Error("[2web] WebSocket error:", { cause: err });
};

console.debug("[2web] Live reload connected");
