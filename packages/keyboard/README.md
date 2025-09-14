# 2Web Kit - Keyboard

```ts
import { ShortcutMap, Shortcut, shift } from "@two-web/kit/keyboard";

function sayHello() {
    console.log("Hello World!");
}

const myShortcuts = new ShortcutMap(
    new Shortcut({
        keys: [shift, "K"],
        action: sayHello,
    }),
);

myShortcuts.listen(document.body);
```
