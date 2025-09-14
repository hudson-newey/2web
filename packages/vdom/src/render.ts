import type { VDomElement } from "./element";

export function render(
  target: HTMLElement,
  ...vdomElements: VDomElement<any, any, any>[]
): void {
  vdomElements.forEach((vdomElement) => {
    const element = vdomElement.toElement();
    target.appendChild(element);
  });
}
