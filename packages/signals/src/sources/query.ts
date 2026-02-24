import { ReadonlySignal } from "../readonlySignal";
import { isSignal } from "../utils/isSignal";
import { unwrapSignal, type MaybeSignal } from "../utils/unwrapSignal";

/**
 * @description
 * A signal that queries the DOM and updates whenever either the query result
 * changes or the queried element changes.
 */
export function query<ElementType extends HTMLElement>(
  querySelector: MaybeSignal<string>
) {
  return new QuerySignal<ElementType>().init(querySelector);
}

class QuerySignal<
  ElementType extends HTMLElement
> extends ReadonlySignal<ElementType | null> {
  private selector: string | null = null;

  public async init(querySelector: MaybeSignal<string>): Promise<this> {
    this.selector = await unwrapSignal(querySelector);

    const initialElement = this.getElement();
    this._internalSet(initialElement);

    const observerConfig = {
      attributes: true,
      attributeFilter: ["class", "style"],
      attributeOldValue: true,
    };

    const fn = async () => {
      this.selector = await unwrapSignal(querySelector);
      this._internalSet(this.getElement());
    };

    const observer = new MutationObserver(fn);
    observer.observe(initialElement ?? document, observerConfig);

    if (isSignal(querySelector)) {
      querySelector.subscribe(fn);
    }

    return this;
  }

  private getElement(): ElementType | null {
    if (!this.selector) {
      return null;
    }

    return document.querySelector<ElementType>(this.selector);
  }
}
