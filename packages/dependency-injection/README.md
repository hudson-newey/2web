# 2Web Kit - Dependency Injection

## Usage

```ts
// main.ts
import { InjectionRoot } from "@two-web/kit/dependency-injection";

const module = new InjectionRoot();
```

```ts
// services/logger.ts
import { provide } from "@two-web/kit/dependency-injection";

export class Logger {
  public constructor(private readonly isProduction = false) {}

  public log(message: string) {
    const logLevel = isProduction ? console.debug : console.log;
    logLevel(message);
  }
}

// This instance will be shared across any consumers.
const sharedInstance = new Logger(true);

provide(Logger, sharedInstance);
```

```ts
// services/fetcher.ts
import { inject } from "@two-web/kit/dependency-injection";
import { Logger } from "@services/logger.ts";

export async function basicFetcher(url: string) {
  const logger = await inject(Logger);

  try {
    return await fetch(url);
  } catch (error) {
    logger.log(error.message);
  }
}
```
