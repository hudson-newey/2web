import type { TwoElement } from "../elements/element";
import { change } from "../renderer/updates";
import type { Directive } from "./directive";

export const when = (
  predicate: boolean | ((...args: any[]) => boolean)
): Directive => {
  return (elementRef: TwoElement) => {
    const predicatePasses =
      typeof predicate === "function" ? predicate() : predicate;

    change(() => {
      if (predicatePasses) {
        elementRef.hidden = undefined;
      } else {
        elementRef.hidden = true;
      }
    });
  };
};
