import type { Character, FunctionType } from "../../typescript";
import { ModifierKey } from "./modifiers";

type ShortcutKey = Character | ModifierKey;

interface IShortcut {
  keys: ShortcutKey[];
  action: FunctionType;

  /**
   * @description
   * Should the `action` callback be invoked repeatedly while the keyboard
   * shortcut is being held down.
   *
   * @default false
   */
  repeat?: boolean;
}

export class Shortcut {
  public readonly keys: Set<ShortcutKey>;
  public readonly action: FunctionType;
  public readonly repeat: boolean;

  /**
   * @description
   * Whether the shortcut is currently being held down.
   */
  public isShortcutHeld = false;

  public constructor(shortcutConfig: IShortcut) {
    this.keys = new Set(shortcutConfig.keys);
    this.action = shortcutConfig.action;
    this.repeat = shortcutConfig.repeat ?? false;
  }

  public matches(event: KeyboardEvent): boolean {
    // Automatically switch ctrl and meta on MacOS.
    // This is because the cmd key is typically used in lieu of ctrl on MacOS.
    // This is because on MacOS, the ctrl key is typically reserved for control
    // characters/commands.
    const isMacOs = navigator.platform.toUpperCase().includes("MAC");
    const isCtrlHeld = isMacOs ? event.metaKey : event.ctrlKey;

    if (this.requiresCtrl() !== isCtrlHeld) return false;
    if (this.requiresShift() !== event.shiftKey) return false;
    if (this.requiresAlt() !== event.altKey) return false;
    if (this.requiresMeta() !== event.metaKey) return false;

    return this.keys.has(event.key as ShortcutKey);
  }

  private requiresCtrl(): boolean {
    return this.keys.has(ModifierKey.Control);
  }

  private requiresShift(): boolean {
    return this.keys.has(ModifierKey.Shift);
  }

  private requiresAlt(): boolean {
    return this.keys.has(ModifierKey.Alt);
  }

  private requiresMeta(): boolean {
    return this.keys.has(ModifierKey.Meta);
  }
}
