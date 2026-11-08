import { expect } from "vitest";
import { Document } from "happy-dom";

export function assertNoErrors(document: Document) {
  const errorElements = document.getElementsByClassName("__2_error_container");
  expect(errorElements.length).toBe(0);
}
