import type { TwoElement } from "../elements/element.ts";
import type { ElementProperty } from "../elements/properties.ts";
import { updateDom } from "../../../_shared/updateDom.ts";
import type { Directive } from "./directive.ts";

type ClassPredicate = () => boolean;
type ClassMap = Readonly<Record<string, ElementProperty | ClassPredicate>>;

export const classMap = (map: ClassMap): Directive => {
  return (elementRef: TwoElement) => {
    Object.entries(map).forEach(([className, propertyOrPredicate]) => {
      const shouldHaveClass =
        typeof propertyOrPredicate === "function"
          ? propertyOrPredicate()
          : Boolean(propertyOrPredicate);

      updateDom(() => {
        if (shouldHaveClass) {
          elementRef.classList.add(className);
        } else {
          elementRef.classList.remove(className);
        }
      });
    });
  };
};
