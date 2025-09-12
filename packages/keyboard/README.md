# 2Web Kit - Keyboard

```ts
import { Keyboard, ShortcutMap, ctrl, shift, optional } from "@two-web/kit/keyboard";

function sayHello() {
    console.log("Hello World!");
}

const myShortcuts = new ShortcutMap({
    [optional(ctrl), shift, "K"]: sayHello,
});

Keyboard.addShortcuts(myShortcuts);
```
