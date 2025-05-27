# 2Web Kit

A collection of packages to enhance your 2web application

| Package      | State |
| ------------ | ----- |
| context      | ❌    |
| database     | ❌    |
| dom-ranges   | ❌    |
| edge         | ❌    |
| ssr          | ❌    |
| pre-fetcher  | ❌    |
| route-guards | ❌    |
| signals      | 🔧    |
| spa-router   | ❌    |
| ssr          | ❌    |
| vite-plugin  | 🔧    |

❌ = Not working, 🔧 = Developer preview, ✅ = Production ready

## Usage

You can install these packages into your project using your package manager.

```sh
$ npm install @two-web/kit
>
```

To use a package from this repository, you can import the sub-package as a an
esm module.

E.g. For JavaScript imports

```js
import { Router } from "@two-web/kit/spa-router";
```

Or for 2web ssg imports

```html
{% include "@two-web/kit/spa-router/spa-router.html" %}
```

Or if you want to use Deno with [JSR](https://jsr.io/)

```js
import { Router } from "https://jsr.io/@two-web/kit/spa-router";
```

## More Information

[GitHub](https://github.com/hudson-newey/2web)
