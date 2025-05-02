import { expect } from "vitest";

export function assertNoErrors(document: Document) {
  const errorElement = document.getElementsByClassName("__2_error_container").length;
  expect(errorElement).toBe(0);
}
