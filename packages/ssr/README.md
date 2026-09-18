# 2Web Kit - SSR

Basic server side rendering (SSR) functionality for 2web/html web applications.

Using SSR will reduce the bundle size of your application, but incur additional
costs to the host.

## Usage

```sh
$ 2web template ssr
>
```

From here, all your applications served through `2web serve` will run through
this ssr server, and `2web build` will produce both a `client/` and `server/`
output.

```ts
import { runServer } from "@two-web/kit/ssr";

runServer({
    port: 2000, // Default: 2000
});
```
