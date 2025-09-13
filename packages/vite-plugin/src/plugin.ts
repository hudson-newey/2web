import { spawn } from "node:child_process";
import { resolve } from "node:path";
import type { Plugin } from "vite";

export function twoWeb(): Plugin {
  return {
    name: "vite-plugin-2web-compiler",
    enforce: "pre",
    async transform(fileContent: string, fileName: string) {
      const binaryPath = resolve(
        process.cwd(),
        "node_modules/two-web/compiler/bin/2webc",
      );
      const args = ["-i", "STDIN", "-o", "STDOUT"];

      const child = spawn(binaryPath, args, {
        stdio: ["pipe", "pipe", "pipe"],
      });

      let stdoutData = "";
      let stderrData = "";

      child.stdout.on("data", (data) => {
        stdoutData += data.toString();
      });

      child.stderr.on("data", (data) => {
        stderrData += data.toString();
      });

      const exitPromise = new Promise<void>((resolve, reject) => {
        child.on("close", (code) => {
          if (code === 0) {
            if (stderrData) {
              console.warn(`2webc stderr for ${fileName}:\n${stderrData}`);
            }
            console.log(
              `2webc finished successfully for ${fileName}. Output size: ${stdoutData.length}`,
            );
            resolve();
          } else {
            console.error(
              `2webc error for ${fileName} (exit code ${code}):\n${stderrData}`,
            );
            reject(
              new Error(
                `2webc failed with exit code ${code}. Stderr:\n${stderrData}`,
              ),
            );
          }
        });

        child.on("error", (err) => {
          console.error(`Failed to start 2webc for ${fileName}:`, err);
          reject(err);
        });
      });

      child.stdin.write(fileContent);
      child.stdin.end();

      try {
        await exitPromise;

        return {
          code: stdoutData,
          map: null,
        };
      } catch (error) {
        throw error;
      }
    },
  };
}
