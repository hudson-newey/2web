import { ReadonlySignal } from "../readonlySignal";

export function form(): FormSignal<any> {
  return new FormSignal({});
}

class FormSignal<T> extends ReadonlySignal<T> {}
