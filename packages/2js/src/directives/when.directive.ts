import type { TwoElement } from "../elements/element";
import { updateDom } from "../../../_shared/updateDom";
import type { Directive } from "./directive";

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
