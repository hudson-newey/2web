import { RouteGuard } from "./predicates";

let routeGuards = new Set<RouteGuard>();

export function addGuard(predicate: RouteGuard): void {
  routeGuards.add(predicate);
}

/**
 * This function is exported so that if you can manually trigger route guards.
 * Note that the guards should automatically trigger on the "load" event.
 *
 * If you are manually calling this function, first check that the order of
 * operations is correct in your application and that this module isn't being
 * dynamically imported.
 */
export async function triggerGuards() {
  for await (const predicate of routeGuards) {
    const predicateOutput = await predicate();

    // We compare the predicate output to "undefined" so that empty strings are
    // a valid redirect back to the home page.
    if (predicateOutput !== undefined) {
      const url = new URL(predicateOutput);

      // We use window.location.replace so that the failed URL doesn't show up
      // in the users "back" history.
      window.location.replace(url.toString());
    }
  }
}

window.addEventListener("load", () => triggerGuards());
