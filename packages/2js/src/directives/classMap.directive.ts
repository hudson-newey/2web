import type { TwoElement } from "../elements/element";
import type { ElementProperty } from "../elements/properties";
import { change } from "../renderer/updates";
import type { Directive } from "./directive";

type ClassPredicate = () => boolean;
type ClassMap = Readonly<Record<string, ElementProperty | ClassPredicate>>;

export const classMap = (map: ClassMap): Directive => {
  return (elementRef: TwoElement) => {
    Object.entries(map).forEach(([className, propertyOrPredicate]) => {
      const shouldHaveClass =
        typeof propertyOrPredicate === "function"
          ? propertyOrPredicate()
          : Boolean(propertyOrPredicate);

      change(() => {
        if (shouldHaveClass) {
          elementRef.classList.add(className);
        } else {
          elementRef.classList.remove(className);
        }
      });
    });
  };
};
