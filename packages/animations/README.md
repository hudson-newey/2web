# 2Web Kit - Animations

## Usage

Vsynced animations that debounces double enqueues.

```ts
import { animate } from "@two-web/kit/animations";

const animationTarget = document.getElementById("ball");

const height = 400;
let progress = 0;
function bounce() {
  const y = Math.abs(Math.sin(progress * Math.PI)) * height;
  animationTarget.style.transform = `translateY(-${y}px)`;

  progress += 0.01;
}

setInterval(() => {
  // Notice that although we perform 16 calls to `animate` per frame
  // (assuming 60fps), the animation callback is only invoked once per
  // frame, because the values are debounced per animation identifier.
  animate(bounce);
}, 1);
```

### Animation Identifiers

Animation identifiers can be used for complex animations that might have
multiple different callers or go across multiple component boundaries.

```ts
import { animation, animate } from "@two-web/kit/animations";

const moveAnimation = animation("move");

function move() {
  // Do something
}

setTimeout(() => {
  animate(move, moveAnimation);
}, 1);
```
