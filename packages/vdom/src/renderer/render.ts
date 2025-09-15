import type { VDomElement } from "../elements/element";

export function render(
  target: HTMLElement,
  ...vdomElements: VDomElement[]
): void {
  vdomElements.forEach((vdomElement) => {
    const element = vdomElement.toElement();
    target.appendChild(element);

    vdomElement.onUpdate();
  });
}
