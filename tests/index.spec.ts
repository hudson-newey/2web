import { test } from './helpers/fixture';
import { navigateToPage } from './helpers/server';

test("should load", async ({ page }) => {
  await navigateToPage(page, "/static-includes.html");
});
