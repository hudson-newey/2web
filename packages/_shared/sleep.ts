/**
 * @summary
 * Returns a promise that is resolved after {@linkcode ms} milliseconds.
 */
export function sleep(ms: number): Promise<void> {
    return new Promise((res) => {
        setTimeout(() => res(), ms);
    });
}
