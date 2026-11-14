import { expect } from "vitest";
import { queryAllByTestId, queryByTestId } from "@testing-library/dom";

const errorContainer = () =>
  queryByTestId(document.body, "__2_error_container");
const errorMessages = () =>
  queryAllByTestId(document.body, "compiler-error-message")
  .map((el) => el.textContent);

export function assertNoErrors() {
  expect(errorContainer()).toBeNull();
}

export function assertHasErrors(expectedErrors: string[]) {
  expect(errorContainer()).toBeDefined();
  expect(errorMessages()).toEqual(expectedErrors);
}
