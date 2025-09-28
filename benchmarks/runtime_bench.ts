#!/usr/bin/env -S deno run --allow-all

// We use tachometer to benchmark the runtime performance of different
// frameworks.
// However, this requires us to inject some code into the frameworks output.
//
// To do this, we Copy the ./node_modules/tachometer/client/lib/bench.js file
// into the dist directory of each framework, then add a script tag to each
// framework's index.html file to load the bench.js file.
// I inject the benchmark code instead of bundling it with the framework so that
// benchmarking code is the same across all frameworks.

// Note that the benchmarks require the tested element to have an id of
// "bench-target".
const injectedCode = `
import * as bench from "/node_modules/tachometer/client/lib/bench.js";

const inTarget = document.getElementById("inTarget");
const outTarget = document.getElementById("outTarget");

let currentCount = 0;
bench.start();
for (let i = 0; i < 100; i++) {
  inTarget.click();

  // Wait until the outTarget updates to the expected value.
  const expectedValue = currentCount + 1;
  while (Number(outTarget.textContent) !== expectedValue) {
    await new Promise((resolve) => requestAnimationFrame(resolve));
  }

  currentCount++;
}
bench.stop();
`;

const frameworks = [
  "2web",
  "preact",
  "react",

  // Svelte is temporarily disabled until I figure out how to make a page build
  // instead of a web component build.
  // "svelte",

  "vanilla",
  "vue",
];

for (const framework of frameworks) {
  const benchPath = `./implementations/${framework}/dist/bench.js`;
  const indexPath = `./implementations/${framework}/dist/index.html`;

  // Copy the bench.js file to the framework's dist directory.
  await Deno.copyFile(
    "./node_modules/tachometer/client/lib/bench.js",
    benchPath
  );

  // Read the framework's index.html file.
  let indexHtml = await Deno.readTextFile(indexPath);

  // Inject the benchmarking code.
  indexHtml = indexHtml.replace(
    "</body>",
    `<script type="module">\n${injectedCode}\n</script>\n</body>`
  );

  // Write the modified index.html file back to the framework's dist directory.
  await Deno.writeTextFile(indexPath, indexHtml);
}
