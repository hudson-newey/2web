const ws = new WebSocket(`ws://${window.location.host}/__2web_updates`);

ws.onmessage = (event) => {
  if (event.data === "reload") {
    console.debug("[2web] Reloading page...");
    window.location.reload();
  }
};

ws.onclose = () => {
  console.error("[2web] Dev server connection lost.");
};

ws.onerror = (err) => {
  throw new Error("[2web] WebSocket error:", { cause: err });
};

console.debug("[2web] Live reload connected");
