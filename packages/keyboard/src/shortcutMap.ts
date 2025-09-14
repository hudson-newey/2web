import type { Shortcut } from "./shortcut";

export class ShortcutMap {
  private readonly shortcuts: Set<Shortcut>;

  public constructor(...shortcuts: Shortcut[]) {
    this.shortcuts = new Set(shortcuts);
  }

  public addShortcut(shortcut: Shortcut) {
    this.shortcuts.add(shortcut);
  }

  public removeShortcut(shortcut: Shortcut) {
    this.shortcuts.delete(shortcut);
  }

  public listen(element: HTMLElement | Document = document) {
    element.addEventListener("keypress", (event) => {
      for (const shortcut of this.shortcuts) {
        if (shortcut.matches(event as KeyboardEvent)) {
          if (shortcut.isShortcutHeld && !shortcut.repeat) return;

          shortcut.isShortcutHeld = true;

          event.preventDefault();
          shortcut.action(event);
        }
      }
    });

    element.addEventListener("keyup", (event) => {
      for (const shortcut of this.shortcuts) {
        if (shortcut.matches(event as KeyboardEvent)) {
          shortcut.isShortcutHeld = false;
        }
      }
    });
  }
}
