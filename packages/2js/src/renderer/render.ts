import type { TwoElement } from "../elements/element";

export function render(
  target: HTMLElement,
  ...twoElements: TwoElement[]
): void {
  twoElements.forEach((twoEl) => {
    const el = twoEl.toElement();
    target.appendChild(el);

    twoEl.onUpdate();
  });
}
