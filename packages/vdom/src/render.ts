import type { VDomElement } from "./element";

export function render(
  target: HTMLElement,
  // biome-ignore lint/suspicious/noExplicitAny: Generic VDomElement type
  ...vdomElements: VDomElement<any, any, any>[]
): void {
  vdomElements.forEach((vdomElement) => {
    const element = vdomElement.toElement();
    target.appendChild(element);
  });
}
