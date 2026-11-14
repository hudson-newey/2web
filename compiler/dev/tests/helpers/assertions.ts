import { expect } from "vitest";
import { queryByTestId } from "@testing-library/dom";

const errorContainer = () =>
  queryByTestId(document.body, "__2_error_container");

export function assertNoErrors() {
  expect(errorContainer()).toBeNull();
}

export function assertHasErrors() {
  expect(errorContainer()).toBeDefined();
}
