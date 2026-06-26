export type ServerEndpoint = string;

export const SERVER_SOCKETS = {
  ACTIONS: "ws://localhost:2000/__2web_actions",
} as const satisfies Record<string, ServerEndpoint>;
