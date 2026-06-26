import * as vscode from "vscode";
import {
  connectDevServer,
  disconnectDevServer,
} from "./dev-server-interop/server-socket";

export function activate(context: Readonly<vscode.ExtensionContext>): void {
  connectDevServer(context);
}

export function deactivate(): void {
  disconnectDevServer();
}
