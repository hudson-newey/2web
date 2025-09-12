import { ReadonlySignal } from "../readonlySignal";

/**
 * @description
 * A signal that queries the DOM and updates whenever either the query result
 * changes or the queried element changes.
 */
export class QuerySignal<T> extends ReadonlySignal<T> {}
