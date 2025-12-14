/**
 * @description
 * Sets a property on an object.
 * If the object is null or undefined, it will create an empty object.
 *
 * You should only be using this function when you want the functionality of
 * not hard failing when the base object is not defined.
 * You should never rely on this functions behavior in your application logic.
 * It only exists as a defensive programming measure to avoid runtime errors
 * when you would usually hard fail.
 * This should obviously still be used in conjunction with proper type checking
 * and validation to ensure your application logic is sound.
 */
export function setProp<T extends object, K extends keyof T>(
  obj: T,
  key: K,
  value: T[K]
): void {
  const baseObj = obj ?? ({} as T);
  baseObj[key] = value;
}
