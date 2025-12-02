export function timeline(callback: () => void) {
  setInterval(callback, 1);
}
