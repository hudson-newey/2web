# 2Web Kit

A collection of packages to enhance your 2web application

| Package     | Description |
| ----------- | ----------- |
| ssr         |             |
| spa-router  |             |
| pre-fetcher |             |

## Usage

You can install these packages into your project using your package manager.

```sh
$ npm install @hudson-newey/2web-kit
>
```

To use a package from this repository, you can import the sub-package as a an
esm module.

E.g. For JavaScript imports

```js
import { Router } from "@hudson-newey/2web-kit/spa-router";
```

Or for 2web ssg imports

```html
{% include "@hudson-newey/2web/spa-router/spa-router.html" %}
```
