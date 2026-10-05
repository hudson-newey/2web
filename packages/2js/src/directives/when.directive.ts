import type { TwoElement } from "../elements/element.ts";
import { updateDom } from "../../../_shared/updateDom.ts";
import type { Directive } from "./directive.ts";

export const when = (
  predicate: boolean | ((...args: any[]) => boolean)
): Directive => {
  return (elementRef: TwoElement) => {
    const predicatePasses =
      typeof predicate === "function" ? predicate() : predicate;

    updateDom(() => {
      if (predicatePasses) {
        elementRef.hidden = undefined;
      } else {
        elementRef.hidden = true;
      }
    });
  };
};
