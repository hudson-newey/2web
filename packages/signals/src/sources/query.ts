import { ReadonlySignal } from "../readonlySignal";

/**
 * @description
 * A signal that queries the DOM and updates whenever either the query result
 * changes or the queried element changes.
 */
export class QuerySignal<
  ElementType extends HTMLElement
> extends ReadonlySignal<ElementType | null> {
  public constructor(private readonly querySelector: string) {
    const initialElement = document.querySelector<ElementType>(querySelector);
    super(initialElement);

    const observerConfig = {
      attributes: true,
      attributeFilter: ["class", "style"],
      attributeOldValue: true,
    };

    const callback = () => {
      this.value = this.getElement();
    };

    const observer = new MutationObserver(callback);
    observer.observe(initialElement ?? document, observerConfig);
  }

  private getElement(): ElementType | null {
    return document.querySelector<ElementType>(this.querySelector);
  }
}
