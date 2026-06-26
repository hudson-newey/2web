export type DevServerCommand = number;

export interface DevServerAction {
  readonly command: Readonly<DevServerCommand>;
}

export const DEV_SERVER_ACTIONS = {
  RECOMPILE_SOURCE: {
    command: 0o1,
  },
  RELOAD_CLIENTS: {
    command: 0o2,
  },
} as const satisfies Record<string, DevServerAction>;
