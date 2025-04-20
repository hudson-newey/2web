import { expect, test as originalTest } from "@playwright/test";

export const test = originalTest.extend({
  page: async ({ page }, use) => {
    // if there is a compiler error, we want the test to fail
    page.on("pageerror", (error) => {
      throw new Error(error.message);
    });

    // similar to page errors, we want to fail if the compiler throws an error
    // this is because 2web injects compiler errors into the document, meaning
    // that the page can still semi-compile with a compiler error, but there
    // will be an injected compiler error.
    const errorElement = await page.locator(".__2_error_container").count();
    expect(errorElement).toBe(0);

    await use(page);
  },
});
