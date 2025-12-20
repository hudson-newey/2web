import { ReadonlySignal } from "../readonlySignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal that queries the DOM and updates whenever either the query result
 * changes or the queried element changes.
 */
export function query<ElementType extends HTMLElement>(
  querySelector: MaybeSignal<string>
) {
  return new QuerySignal<ElementType>(querySelector);
}

class QuerySignal<
  ElementType extends HTMLElement
> extends ReadonlySignal<ElementType | null> {
  private readonly selector: string;

  public constructor(querySelector: MaybeSignal<string>) {
    const initialElement = document.querySelector<ElementType>(
      unwrapSignal(querySelector)
    );

    super(initialElement);

    this.selector = unwrapSignal(querySelector);
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
    return document.querySelector<ElementType>(this.selector);
  }
}
