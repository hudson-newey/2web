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
import * as bench from "/bench.js";

// Because some frameworks (React, Vue, etc...) bootstrap themselves using
// JavaScript, we need to wait until the DOM has been fully populated with our
// target elements before running the benchmark.
let inTarget;
let outTarget;
do {
  await new Promise((resolve) => setTimeout(resolve, 0));
  inTarget = document.getElementById("inTarget");
  outTarget = document.getElementById("outTarget");
} while (!inTarget || !outTarget);

let currentCount = 0;
bench.start();
for (let i = 0; i < 100; i++) {
  inTarget.click();

  // Wait until the outTarget updates to the expected value.
  const expectedValue = currentCount + 1;
  while (Number(outTarget.textContent) !== expectedValue) {
    await new Promise((resolve) => setTimeout(resolve, 0));
  }

  currentCount++;
}
bench.stop();
`;

const frameworks = [
  "2web",
  "preact",
  "react",
  "svelte",
  "vanilla",
  "vue",
];

// Create the dist_bench directory if it doesn't exist. If  it already exists,
// delete it first
try {
  await Deno.remove("./dist_bench", { recursive: true });
}
catch {
  // do nothing
}

await Deno.mkdir("./dist_bench", { recursive: true });

// Copy the bench.js file to the framework's dist directory.
await Deno.copyFile(
  "./node_modules/tachometer/client/lib/bench.js",
  "./dist_bench/bench.js"
);

function copyDirectory(src: string, dest: string) {
  return Deno.mkdir(dest, { recursive: true }).then(async () => {
    for await (const entry of Deno.readDir(src)) {
      const srcPath = `${src}/${entry.name}`;
      const destPath = `${dest}/${entry.name}`;

      if (entry.isFile) {
        await Deno.copyFile(srcPath, destPath);
      } else if (entry.isDirectory) {
        await copyDirectory(srcPath, destPath);
      }
    }
  });
}

for (const framework of frameworks) {
  const distPath = `./implementations/${framework}/dist`;
  const indexPath = `${distPath}/index.html`;

  // Read the framework's index.html file.
  let indexHtml = await Deno.readTextFile(indexPath);

  // Inject the benchmarking code.
  indexHtml = indexHtml.replace(
    "</body>",
    `<script type="module">\n${injectedCode}\n</script>\n</body>`
  );

  // Write the modified index.html file back to the framework's dist directory.
  await Deno.writeTextFile(indexPath, indexHtml);

  // Copy all files from the framework's dist directory to a new directory
  // called dist_bench.
  const distBenchPath = `./dist_bench/`;
  await Deno.mkdir(distBenchPath, { recursive: true });

  for await (const entry of Deno.readDir(distPath)) {
    const srcPath = `${distPath}/${entry.name}`;

    // If we are copying the index.html file, replace it with the framework name
    const destPath = entry.name === "index.html"
      ? `${distBenchPath}/${framework}.html`
      : `${distBenchPath}/${entry.name}`;

    if (entry.isFile) {
      await Deno.copyFile(srcPath, destPath);
    } else if (entry.isDirectory) {
      await copyDirectory(srcPath, destPath);
    }
  }
}
