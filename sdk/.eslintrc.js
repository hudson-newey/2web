import { defineConfig } from "eslint/config";

export default defineConfig([
  {
    files: ["**/*.js", "**/*.ts"],
    extends: ["js/recommended"],
    rules: {
			"no-console": "error",
      "camelcase": "error",
      "capitalized-comments": "error",
      "no-eval": "error",

      "no-unused-vars": "warn",
    },
  },
]);
