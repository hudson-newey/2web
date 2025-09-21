# 2Web Kit - Hydration

SSR partial hydration primitives.

This feature is not available in web workers.
We do not support transferring web worker state across ssr boundaries.

## Usage

```ts
import { hydrate } from "@two-web/kit/hydration";

const params = new URLSearchParams(window.location.search);
const userId= params.get(`user`);

// Using the "hydrate" decorator means that this variable will be fetched on the
// server and then hydrated on the client, meaning that we don't need to fetch
// it client side unless SSR times out.
@hydrate()
const userModel = await fetch(`https://api.example.com/users/${id}`, {
    method: "POST",
    headers: { "Accept": "application/json" },
  })
  .then((resp) => resp.json());
```
