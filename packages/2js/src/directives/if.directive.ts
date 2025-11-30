import type { TwoElement } from "../elements/element";
import type { Directive } from "./directive";

export const iif = (
  predicate: boolean | ((...args: any[]) => boolean)
): Directive => {
  return (elementRef: TwoElement) => {
    const predicatePasses =
      typeof predicate === "function" ? predicate() : predicate;

    if (predicatePasses) {
      elementRef.hidden = undefined;
    } else {
      elementRef.hidden = true;
    }
  };
};
