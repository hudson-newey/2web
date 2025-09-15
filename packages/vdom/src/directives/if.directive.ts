import type { VDomElement } from "../elements/element";
import type { Directive } from "./directive";

export const vDomIf = (predicate: boolean | ((...args: any[]) => boolean)): Directive  => {
  return (vDomElement: VDomElement) => {
    const predicatePasses = typeof predicate === "function" ? predicate() : predicate;

    if (predicatePasses) {
      vDomElement.hidden = undefined;
    } else {
      vDomElement.hidden = true;
    }
  };
};
