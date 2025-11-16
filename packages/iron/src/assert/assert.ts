export function assert(assertion: () => boolean, message?: string): void {
  if (!assertion()) {
    throw new Error(message || "Assertion failed");
  }
}
