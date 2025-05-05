type FrameworkName = string;
type Kilobyte = number;

const reportsDirectory = "reports/";

function implementationPath<T extends FrameworkName>(name: T) {
  return `implementations/${name}` as const;
}

async function buildImplementation(name: string) {
  const command = new Deno.Command("bash", {
    args: [`${implementationPath(name)}/build.sh`],
  });

  const { code, stdout, stderr } = await command.output();
  // if (code !== 0) {
  //   console.log("Error:", code);
  //   throw new Error(String.fromCharCode(...stderr));
  // }

  console.debug(String.fromCharCode(...stdout));
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
    for await (const entry of Deno.readDir(dirPath)) {
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
  const testedFrameworks = [
    "2web",
    "svelte",
    "preact",
    "vue",
    "react",
  ] as const satisfies FrameworkName[];

  const results = Promise.all(
    testedFrameworks.map((name) => buildImplementation(name)),
  );

  await results;

  // we delete the original reports/ directory and create a new one because
  // I didn't want to code up partial report generation
  try {
    await Deno.lstat(reportsDirectory);

    // reports directory exists
    await Deno.remove(reportsDirectory, { recursive: true });
  } catch {
    // do nothing
  }

  await Deno.mkdir(reportsDirectory, { recursive: true });

  // I write to a csv so that I can do some cool graphs in the future and so
  // that loading into excel is easy
  const assetSizeFileName = reportsDirectory + "assetSize.csv";
  await Deno.writeTextFile(assetSizeFileName, "Framework,size (KB)\n", {
    append: true,
  });

  // we do this in async because I'm too tired to code it otherwise
  for (const implementation of testedFrameworks) {
    const distPath = implementationPath(implementation) + "/dist/";
    const distSize = await getDirectorySize(distPath);

    await Deno.writeTextFile(
      assetSizeFileName,
      `${implementation},${distSize}\n`,
      { append: true },
    );
  }

  console.log("Benchmark reports successfully generated!");
}

await runBenchmark();
