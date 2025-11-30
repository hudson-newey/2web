// In the constructor we use a plain object for easier construction

import type { Directive } from "../directives/directive";
import { change } from "../renderer/updates";

// but internally we use Maps for easier updates.
interface TwoElementConstructor {
  tagName: string;
  textContent?: string;
  attributes?: Record<string, string>;
  properties?: Record<string | symbol, unknown>;
  events?: Record<string, (event: Event) => void>;

  directives?: Directive[];

  children?: TwoElement[];
}

interface FreeformElement {
  [key: string]: unknown;
}

export class TwoElement implements FreeformElement {
  private readonly tagName: string;
  private readonly attributes: Map<string, string>;
  private readonly properties: Map<string | symbol, unknown>;
  private readonly events: Map<string, (event: Event) => void>;
  private readonly directives: Directive[];
  private readonly children: TwoElement[];
  private readonly proxy: typeof this;
  [key: string]: unknown;

  private textContent: string;
  private ref: HTMLElement | null = null;

  public constructor({
    tagName,
    textContent = "",
    attributes = {},
    properties = {},
    events = {},
    directives = [],
    children = [],
  }: TwoElementConstructor) {
    this.tagName = tagName;
    this.textContent = textContent;

    // Convert the plain objects to Maps for easier updates
    this.attributes = new Map(Object.entries(attributes));
    this.properties = new Map(Object.entries(properties));
    this.events = new Map(Object.entries(events));
    this.directives = directives;
    this.children = children;

    this.proxy = new Proxy(this, {
      set: (target, property, value) => {
        if (property === "ref") {
          Reflect.set(target, property, value);
          return true;
        }

        if (property === "textContent") {
          target.setTextContent(value as string);
          return true;
        }

        if (typeof property === "string") {
          if (this.attributes.has(property)) {
            this.setAttribute(property, value);
            return true;
          }
        }

        this.setProperty(property, value);
        return true;
      },
    });

    return this.proxy;
  }

  // TODO: This method should return an promise that resolves after the first
  // update cycle is complete, and the element is attached to the DOM.
  /**
   * @returns
   * A read-only HTMLElement constructed from the TwoElement's properties.
   * This element is not attached to the DOM.
   * We use a readonly type so that consumers don't accidentally modify 2js
   * elements 2js outside of the 2js system.
   *
   * If you need to modify this element, you should not be modifying the
   * HTMLElement's properties or attributes directly, but instead modify the
   * TwoElement's virtual properties, attributes, event listeners, etc...
   */
  public toElement(): Readonly<HTMLElement> {
    const element = document.createElement(this.tagName);

    // Note that this ref assignment occurs before any DOM updates occur.
    // Meaning that if you want to get the initial state of the live DOM, you
    // will need to wait until after the first update cycle.
    //
    // This should not be a problem however, because we support sub-frame
    // updates if you directly read/write to the TwoElement properties.
    this.ref = element;

    change(() => {
      element.textContent = this.textContent;
      this.attributes.forEach((value, key) => {
        element.setAttribute(key, value);
      });

      this.properties.forEach((value, key) => {
        (element as any)[key] = value;
      });

      this.events.forEach((handler, eventName) => {
        element.addEventListener(eventName, handler);
      });

      for (const child of this.children) {
        element.appendChild(child.toElement());
      }
    });

    return Object.freeze(this.ref);
  }

  public onUpdate() {
    this.directives.forEach((directive) => {
      directive(this.proxy);
    });
  }

  // Note that in all of these "update" methods, we update the actual DOM first
  // so if updating the DOM throws an error, we don't update the TwoElement
  // state.

  private setProperty(key: string | symbol, value: unknown) {
    change(() => {
      if (this.ref) {
        (this.ref as any)[key] = value;
      }
    });

    const currentValue = this.properties.get(key);
    if (currentValue === value) {
      return;
    }

    // Note that we update the internal state before the actual DOM is updated
    // meaning that we can perform sub-frame updates that are dependent on each
    // other.
    this.properties.set(key, value);
    this.onUpdate();
  }

  private setAttribute(key: string, value: string) {
    change(() => {
      if (this.ref) {
        this.ref.setAttribute(key, value);
      }
    });

    this.attributes.set(key, value);
    this.onUpdate();
  }

  private setTextContent(value: string) {
    change(() => {
      if (this.ref) {
        this.ref.textContent = value;
      }
    });

    this.textContent = value;
    this.onUpdate();
  }
}
