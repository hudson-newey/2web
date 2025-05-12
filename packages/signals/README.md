# 2Web Signals

Lightweight framework agnostic signals that provides **runtime** state tracking.

## Usage

```html
<script>
import { Signal, effect } from "@two-web/kit/signals";

const count = new Signal(0);

effect(() => {
  console.log("New value is ${count.value}");
}, [count]);
</script>

<button onclick="count.set(count.value + 1)">Increment</button>
<button onclick="count.set(count.value - 1)">Decrement</button>
```
