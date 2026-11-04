#!/usr/bin/env -S deno run --allow-run --allow-write --allow-read
import { FrameworkName, getFrameworks } from "../src//constants/frameworks.ts";
import { ASSETS } from "../src/constants/assets.ts";

type Kilobyte = number;
type Milliseconds = number;

const reportsDirectory = ASSETS.outDir + "reports/";

function implementationPath<T extends FrameworkName>(name: T) {
  return `implementations/${name}` as const;
}

/**
 * @returns
 * The amount of time (in milliseconds) that the build took.
 * Limitation: This includes the time that it takes Deno to call bash, run the
 * script, etc...
 * Build times are therefore not very scientific, but a good smoke test to
 * see if there are any build time outliers.
 */
async function buildImplementation(name: string): Promise<Milliseconds> {
  const t0 = performance.now();
  const command = new Deno.Command("bash", {
    args: [`${implementationPath(name)}/build.sh`],
  });

  const { stdout } = await command.output();
  const t1 = performance.now();

  console.debug(String.fromCharCode(...stdout));

  // see the jsdoc above This is no where near a good benchmark, but it exists
  // to catch outliers in build time.
  return t1 - t0;
}

// warning: this test is "bad" because it doesn't track the bundle size but
// instead tracks the size of all built assets
//
// This is also AI generated code which is why the code quality is so awful
//
// TODO: I should improve this benchmark test
async function getDirectorySize(dirPath: string): Promise<Kilobyte> {
  let totalSize = 0;

  try {
    const dir = Deno.readDirSync(dirPath);
    for (const entry of dir) {
      const path = `${dirPath}/${entry.name}`;

      if (entry.isFile) {
        const stat = await Deno.stat(path);

        // divide by 1000 here because we want to report the size back in KB
        // not the Deno default Bytes
        totalSize += stat.size / 1_000;
      } else if (entry.isDirectory) {
        totalSize += await getDirectorySize(path);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dirPath}:`, error);
    throw error;
  }

  return totalSize;
}

async function runBenchmark() {
  const targets = getFrameworks();

  // Delete the dist/ directory if it exists
  for (const framework of targets) {
    const distPath = implementationPath(framework) + "/dist";

    try {
      Deno.removeSync(distPath, { recursive: true });
    } catch {
      // do nothing
    }
  }

  const buildOutputs = await Promise.all(
    targets.map(async (name) => {
      const buildTime = await buildImplementation(name);
      return { name, buildTime };
    }),
  );

  const buildTimes: Partial<Record<string, Milliseconds>> = {};
  for (const build of buildOutputs) {
    buildTimes[build.name] = build.buildTime;
  }

  // we delete the original reports/ directory and create a new one because
  // I didn't want to code up partial report generation
  try {
    await Deno.lstat(reportsDirectory);

    // reports directory exists
    Deno.removeSync(reportsDirectory, { recursive: true });
  } catch {
    // do nothing
  }

  Deno.mkdirSync(reportsDirectory, { recursive: true });

  // I write to a csv so that I can do some cool graphs in the future and so
  // that loading into excel is easy
  const assetSizeFileName = reportsDirectory + "build_report.csv";
  Deno.writeTextFileSync(
    assetSizeFileName,
    "Framework,size (KB),Build Time (ms)\n",
    {
      append: true,
    },
  );

  // we do this in async because I'm too tired to code it otherwise
  for (const implementation of targets) {
    const distPath = implementationPath(implementation) + "/dist/";
    const distSize = await getDirectorySize(distPath);
    const buildTime = buildTimes[implementation];

    Deno.writeTextFileSync(
      assetSizeFileName,
      `${implementation},${distSize},${buildTime}\n`,
      { append: true },
    );
  }

  console.log("Benchmark reports successfully generated!");
}

await runBenchmark();
