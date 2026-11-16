import { beforeEach, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";
import { userEvent } from "@testing-library/user-event";

const countBtn = () => getByTestId(document.body, "count-btn");

beforeEach(async () => {
  await navigateToPage("textnode-reactivity.html");
});

it("has no errors", () => {
  assertNoErrors();
});

it("has the correct initial value", () => {
  expect(countBtn()).toHaveTextContent("0");
});

it("increases when clicked", async () => {
  const user = userEvent.setup({ document });

  await user.click(countBtn());
  expect(countBtn()).toHaveTextContent("1");

  // Make sure it's not a fluke that worked once :P
  await user.click(countBtn());
  expect(countBtn()).toHaveTextContent("2");
});
