import type { TwoElement } from "../elements/element.ts";
import type { Directive } from "./directive.ts";

type StyleName = string;
type StyleValue = string;
type StyleCallback = () => string;
type StyleMap = Readonly<Record<StyleName, StyleValue | StyleCallback>>;

export const styleMap = (map: StyleMap): Directive => {
  return (elementRef: TwoElement) => {
    Object.entries(map).forEach(([key, value]) => {
      const styleValue = typeof value === "function" ? value() : value;
      elementRef.style[key] = styleValue;
    });
  };
};
