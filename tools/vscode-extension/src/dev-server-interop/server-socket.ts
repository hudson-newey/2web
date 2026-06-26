import * as vscode from "vscode";
import WebSocket from "ws";
import { SERVER_SOCKETS } from "./server-endpoints";
import { DEV_SERVER_ACTIONS, DevServerAction } from "./server-commands";

let ws: WebSocket | null = null;

export function connectDevServer(
  context: Readonly<vscode.ExtensionContext>,
): void {
  connectWebSocket();

  const supportedLanguages = ["2web", "html", "typescript", "javascript"];

  const disposable = vscode.workspace.onDidSaveTextDocument(
    (document: vscode.TextDocument) => {
      if (supportedLanguages.includes(document.languageId)) {
        dispatchAction(DEV_SERVER_ACTIONS.RECOMPILE_SOURCE);
      }
    },
  );

  context.subscriptions.push(disposable);
}

export function disconnectDevServer(): void {
  if (ws) {
    ws.close();
  }
}

function connectWebSocket(): void {
  ws = new WebSocket(SERVER_SOCKETS.ACTIONS);

  ws.on("open", () => {
    console.log("WebSocket connection established.");
  });

  ws.on("close", () => {
    console.log("WebSocket connection closed. Retrying in 5 seconds...");
    setTimeout(connectWebSocket, 5000);
  });

  ws.on("error", (error: Error) => {
    console.error("WebSocket error:", error.message);
  });
}

function dispatchAction(action: Readonly<DevServerAction>): void {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;

  ws.send(action.command, (error) => {
    console.error("failed to dispatch dev server action", { action, error });
  });
}
