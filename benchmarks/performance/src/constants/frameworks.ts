export function getFrameworks(): string[] {
  const implementationPath = "implementations/";

  const entries = Deno.readDirSync(implementationPath);
  return [...entries].filter((x) => x.isDirectory).map((x) => x.name);
}
