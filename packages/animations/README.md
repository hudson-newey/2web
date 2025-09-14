# 2Web Kit - Animations

Vsynced animations that debounces double enqueues.

```ts
import { animation, animate } from "@two-web/kit/animations";

const bounce = animation("bounce");

const animationTarget = document.getElementById("ball");

let progress = 0;
setInterval(() => {
  // Notice that although we perform 16 calls to `animate` per frame
  // (assuming 60fps), the animation callback is only invoked once per
  // frame, because the values are debounced per animation identifier.
  animate(bounce, () => {
    const height = 400;
    const y = Math.abs(Math.sin(progress * Math.PI)) * height;
    animationTarget.style.transform = `translateY(-${y}px)`;

    progress += 0.01;
  });
}, 1);
```
