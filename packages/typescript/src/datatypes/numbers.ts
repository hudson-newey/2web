import { Not } from "../conditions";

/**
 * @description
 * Removes "Infinity", "NaN"s from the number type.
 */
export type ValidNumber = number & object;

const a: ValidNumber = 42;
const b: ValidNumber = Infinity;
