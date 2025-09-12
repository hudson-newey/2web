# 2Web Kit

A collection of packages to enhance your 2web application

| Package                                        | State |
| ---------------------------------------------- | ----- |
| [animations](animations/README.md)             | 🔧     |
| [browser-state](event-listener/README.md)      | ❌     |
| [event-listener](event-listener/README.md)     | ❌     |
| [keyboard](keyboard/README.md)                 | ❌     |
| [pre-fetcher](pre-fetcher/README.md)           | 🔧     |
| [route-guards](route-guards/README.md)         | 🔧     |
| [signals](signals/README.md)                   | 🔧     |
| [ssr](ssr/README.md)                           | 🔧     |
| [threads](threads/README.md)                   | ❌     |
| [typescript](threads/README.md)                | 🔧     |
| [view-transitions](view-transitions/README.md) | ❌     |
| [vite-plugin](vite-plugin/README.md)           | 🔧     |

❌ = Not working, 🔧 = Developer preview, ✅ = Production ready

## Usage

### Package Managers

You can install these packages into your project using your package manager.

```sh
$ npm install @two-web/kit
>
```

**OR** you can add the following tag to your page's `<head>` tag.

```html
<script type="module" src="https://cdn.jsdelivr.net/npm/@two-web/kit"></script>
```

To use a package from this repository (either through node_modules or CDN), you
can import the sub-package as a an esm module.

```html
<script type="module">
import { EventHandler } from "@two-web/kit/signals";

const target = document.getElementById("counter");
const countHandler = new EventHandler((event, value) => {
  const count = value + 1;

  event.target.textContent = count;
  return count;
});

target.addEventListener("click", countHandler);
</script>

<button id="counter">0</button>
```
