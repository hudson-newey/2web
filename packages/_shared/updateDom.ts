type Change = () => void;
type RafRef = ReturnType<typeof requestAnimationFrame>;

const updateQueue = new Set<Change>();
let queuedUpdate: RafRef | null = null;

/**
 * @description
 * An optimized update function that batches multiple updates into a single
 * animation frame.
 *
 * @param mod - The update callback to be executed.
 */
export function updateDom(mod: Change): RafRef {
  updateQueue.add(mod);

  // We only enqueue an update frame if one isn't already queued.
  // If there is already once queued, the update will be processed in the next
  // queued frame.
  queuedUpdate ??= requestAnimationFrame(() => {
    runUpdate();
  });

  return queuedUpdate;
}

function runUpdate(): void {
  for (const update of updateQueue) {
    try {
      update();
    } catch (cause) {
      // Use console.error instead of throwing so that we don't cancel the rest
      // of the updates in the queue.
      console.error(cause);
    }
  }

  // We perform one big cleanup at the end instead of removing items as we go to
  // avoid mutating the set while iterating over it and to reduce the number of
  // operations.
  // Although this does mean that if an update fails, the failed update will be
  // removed from the queue and not retried.
  updateQueue.clear();
  queuedUpdate = null;
}
