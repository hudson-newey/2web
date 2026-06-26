export type DevServerCommand = number;

export interface DevServerAction {
  readonly command: Readonly<DevServerCommand>;
}

export const DEV_SERVER_ACTIONS = {
  RELOAD_CLIENTS: {
    command: 0o1,
  },
} as const satisfies Record<string, DevServerAction>;
