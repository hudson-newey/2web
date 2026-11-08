import { expect } from "vitest";
import { screen } from "@testing-library/dom";

export function assertNoErrors() {
  const errorContainer = screen.queryByTestId("__2_error_container");
  expect(errorContainer).toBeNull();
}
