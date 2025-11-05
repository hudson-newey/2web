# TypeScript

Similar to most modern frameworks, we support creating websites with TypeScript.
By default, all `<script>` tags are TypeScript-enabled, meaning that you don't
have to set up a build pipeline, set up any configuration, or worry about
increasing the complexity.

Since TypeScript is a superset of JavaScript, if you don't want to write
TypeScript code, you can simply not.
If you want to use [JsDoc](https://jsdoc.app/) types, you can as well.
By default, JsDoc comments will be removed when the `--production` flag is
used.

By default, we use [esbuild](https://esbuild.github.io/) for lightning fast
TypeScript builds, meaning that we don't provide any type checking.
Type checking is considered a development environment operation, and should not
be performed at compile time.

By default, all `<script>` tags are compiled to
[ecma script modules](https://webpack.js.org/guides/ecma-script-modules/) (esm).
This means that you can import any third party libraries that export JavaScript
or TypeScript as ESM.

## 2Web Kit Helpers

In the [`@two-web/kit/typescript`](../packages/typescript/README.md) package,
there exists typescript helpers for common use cases such as untrusted inputs,
unstable variables, and testing types.

[**Next**](./9-styling.md) (Page/Component Styling)

## Route scripts

If you want to run a script in all pages within a directory (not including
sub-directories), you can add a `__scripts.ts` file to the containing directory.

```txt
└─ .
   ├─ __scripts.ts
   └─ article.html
```
