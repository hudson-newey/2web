import type { AnimationIdentifier } from "./animation";

export type AnimationCallback = (...args: any[]) => unknown;

export function animate(
  identifier: AnimationIdentifier,
  callback: AnimationCallback,
) {
  vsyncQueue.set(identifier, callback);
  frame();
}

// TODO: Replace this with a WeakMap
const vsyncQueue = new Map<AnimationIdentifier, AnimationCallback>();

let hasQueuedFrame = false;
function frame() {
  if (hasQueuedFrame) {
    return;
  }

  hasQueuedFrame = true;

  requestAnimationFrame(() => {
    // We set the queued frame reference so that if the current frame takes
    // a long time to process, we don't drop a frame / state.
    hasQueuedFrame = false;

    for (const [identifier, callback] of vsyncQueue.entries()) {
      callback();
      vsyncQueue.delete(identifier);
    }
  });
}
