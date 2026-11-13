import { beforeEach, describe, expect, it } from "vitest";
import { assertNoErrors } from "./tests/helpers/assertions";
import { navigateToPage } from "./tests/helpers/fixture";
import { getByTestId } from "@testing-library/dom";

beforeEach(async () => {
  await navigateToPage("code-blocks.html");
});

it("has no errors", () => {
  assertNoErrors();
});

describe("empty code blocks", () => {
  it("renders an empty code block", () => {
    const el = getByTestId(document.body, "empty-code");
    expect(el.innerHTML).toBe("");
  });

  it("renders an empty script", () => {
    const el = getByTestId(document.body, "empty-script");
    expect(el.innerHTML).toBe("&gt;script&lt;&gt;/script&lt;");
  });
});

describe("css code blocks", () => {
  it("escapes code", () => {
    const cssCodeBlock = getByTestId(document.body, "css-code");
    expect(cssCodeBlock.innerHTML).toBe(`<code>
  &lt;style&gt;
    .unstyledElement {
      color: red !important;
    }
  &lt;/style&gt;
</code>`);
  });

  it("doesn't apply styling to targeted elements", () => {
    const effectedElement = getByTestId(document.body, "unstyled-element");
    expect(effectedElement.style.color).toBe("green");
  });
});
