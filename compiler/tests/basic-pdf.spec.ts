import { readFileSync } from "node:fs";
import { expect, test } from "vitest";

// Binary assets are passed through to the build output byte for byte. The pdf
// can't be loaded through the html test harness, so the passthrough is
// asserted directly against the emitted file.
test("should pass the pdf through unmodified", () => {
  const compiled = readFileSync("dist/basic-pdf.pdf");
  const source = readFileSync("dev/basic-pdf.pdf");

  expect(compiled.equals(source)).toBe(true);
});

test("should keep the pdf file signature", () => {
  const compiled = readFileSync("dist/basic-pdf.pdf");

  expect(compiled.subarray(0, 5).toString("ascii")).toEqual("%PDF-");
});
