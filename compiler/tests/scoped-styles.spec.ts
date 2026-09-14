import { readFileSync } from "node:fs";
import { beforeEach, expect, test } from "vitest";
import { navigateToPage } from "./helpers/fixture";
import { assertNoErrors } from "./helpers/assertions";
import { Document } from "happy-dom";

// Component styles are scoped with the css @scope at rule: every instance of
// a component gets a unique scoping attribute on its top level elements, and
// the component's style blocks are wrapped in "@scope" rules for that
// attribute, so the component's styles can't leak into the rest of the page.

let document: Document;

beforeEach(async () => {
  document = (await navigateToPage("virtual-functions.html")).document;
});

test("should load", () => {
  assertNoErrors(document);
});

test("should tag component instances with their scope attribute", () => {
  const badges = document.querySelectorAll(".badge-label");

  expect(badges.length).toBe(2);
  expect(badges[0].getAttribute("data-__2_scope")).toBe("1");
  expect(badges[1].getAttribute("data-__2_scope")).toBe("2");
});

test("should not tag the page's own elements with a scope attribute", () => {
  expect(document.querySelector(".page-paragraph")!.getAttribute("data-__2_scope")).toBeNull();
  expect(document.querySelector(".ids")!.getAttribute("data-__2_scope")).toBeNull();
});

test("should wrap the component's style block in a @scope rule per instance", () => {
  // The css assets that the page references contain the scoped rules.
  const css = [...document.querySelectorAll('link[rel="stylesheet"]')]
    .map((link) => link.getAttribute("href")!)
    .map((href) => readFileSync(`dist/${href}`, "utf-8"))
    .join("\n");

  expect(css).toContain('@scope ([data-__2_scope="1"])');
  expect(css).toContain('@scope ([data-__2_scope="2"])');
  expect(css.match(/@scope/g)?.length).toBe(2);

  // The scoped rules are the component's own rules.
  expect(css).toContain("background-color: darkslateblue");
});
