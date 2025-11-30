type Change = () => void;

const updateQueue = new Set<Change>();
let queuedUpdate: ReturnType<typeof requestAnimationFrame> | null = null;

/**
 * @description
 * An optimized update function that batches multiple updates into a single
 * animation frame.
 *
 * @param mod - The update callback to be executed.
 */
export function change(mod: Change): void {
  updateQueue.add(mod);

  // We only enqueue an update frame if one isn't already queued.
  // If there is already once queued, the update will be processed in the next
  // queued frame.
  queuedUpdate ??= requestAnimationFrame(() => {
    runUpdate();
  });
}

function runUpdate(): void {
  updateQueue.forEach((update) => {
    try {
      update();
    } catch (cause) {
      // Use console.error instead of throwing so that we don't cancel the rest
      // of the updates in the queue.
      console.error(cause);
    }
  });

  updateQueue.clear();
  queuedUpdate = null;
}
