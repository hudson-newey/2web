// In the constructor we use a plain object for easier construction
// but internally we use Maps for easier updates.
interface VDomElementConstructor {
  tagName: string;
  textContent?: string;
  attributes?: Record<string, string>;
  properties?: Record<string | symbol, unknown>;
  events?: Record<string, (event: Event) => void>;
}

export class VDomElement {
  private readonly tagName: string;
  private readonly attributes: Map<string, string>;
  private readonly properties: Map<string | symbol, unknown>;
  private readonly events: Map<string, (event: Event) => void>;

  private textContent: string;
  private ref: HTMLElement | null = null;

  public constructor(
    {
      tagName,
      textContent = "",
      attributes = {},
      properties = {},
      events = {},
    } = {} as VDomElementConstructor,
  ) {
    this.tagName = tagName;
    this.textContent = textContent;

    // Convert the plain objects to Maps for easier updates
    this.attributes = new Map(Object.entries(attributes));
    this.properties = new Map(Object.entries(properties));
    this.events = new Map(Object.entries(events));

    // biome-ignore lint/correctness/noConstructorReturn: Proxy-based reactivity is expected
    return new Proxy(this, {
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
  }

  /**
   * @returns
   * A read-only HTMLElement constructed from the VDomElement's properties.
   * This element is not attached to the DOM.
   * We use a readonly type so that consumers don't accidentally modify vdom
   * elements VDom outside of the VDom system.
   */
  public toElement(): Readonly<HTMLElement> {
    const element = document.createElement(this.tagName);

    element.textContent = this.textContent;
    this.attributes.forEach((value, key) => {
      element.setAttribute(key, value);
    });

    this.properties.forEach((value, key) => {
      // biome-ignore lint/suspicious/noExplicitAny: TODO: Better typing for properties
      (element as any)[key] = value;
    });

    this.events.forEach((handler, eventName) => {
      element.addEventListener(eventName, handler);
    });

    this.ref = element;
    return Object.freeze(this.ref);
  }

  // Note that in all of these "update" methods, we update the actual DOM first
  // so if updating the DOM throws an error, we don't update the VDomElement
  // state.

  private setProperty(key: string | symbol, value: unknown) {
    if (this.ref) {
      // biome-ignore lint/suspicious/noExplicitAny: TODO: Better typing for properties
      (this.ref as any)[key] = value;
    }

    this.properties.set(key, value);
  }

  private setAttribute(key: string, value: string) {
    if (this.ref) {
      this.ref.setAttribute(key, value);
    }

    this.attributes.set(key, value);
  }

  private setTextContent(value: string) {
    if (this.ref) {
      this.ref.textContent = value;
    }
    this.textContent = value;
  }
}
