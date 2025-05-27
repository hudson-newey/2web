type RedirectLocation = Parameters<Location["replace"]>[0];
type MaybeAsync<T> = T | Promise<T>;

/**
 * @returns undefined - If the route guard was successful
 * @returns RedirectLocation - If the user needs to redirect to another page
 */
export type RouteGuard = () =>
  | MaybeAsync<undefined>
  | MaybeAsync<RedirectLocation>;
