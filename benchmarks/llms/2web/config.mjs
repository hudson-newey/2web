import { getBuiltInRatings } from "web-codegen-scorer";

/** @type {import("web-codegen-scorer").EnvironmentConfig} */
export default {
  displayName: "2web",
  clientSideFramework: "2web",
  sourceDirectory: "",
  ratings: [
    // This includes some framework-agnostic scoring ratings to your eval.
    // You can add your own custom ratings to this array.
    ...getBuiltInRatings(),
  ],
  generationSystemPrompt: "./system-instructions.md",
  executablePrompts: ["./prompts/**/*.md"],

  // The following options aren't mandatory, but can be useful:
  // id: '', Unique ID for the environment. If empty, one is generated from the `displayName`.
  // packageManager: 'npm', // Name of the package manager used to install dependencies.
  // skipInstall: false, // Whether to skip installing dependencies. Useful if you're doing it yourself already.
  // buildCommand: 'npm run build', // Command used to build the generated code.
  // serveCommand: 'npm run start -- --port 0', // Command used to start a dev server with the generated code.
  // mcpServers: [], // Model Context Protocal servers to run during the eval.

  // repairSystemPrompt: '', // Path to a prompt used when repairing broken code.
  // editingSystemPrompt: '', // Path to a prompt used when editing code during a multi-step eval.
  // codeRatingPrompt: '', // Path to a prompt to use when automatically rating the generated code with an LLM.
  // classifyPrompts: false, // Whether to exclude the prompt text from the final report.

  // Path to a directory that will be merged with the `sourceDirectory` to produce
  // the final project. Useful for reusing boilerplate between environments.
  // projectTemplate: '',

  // If your setup has different client-side and full-stack framework,
  // you can specify a different full-stack framework here.
  // fullStackFramework: '',
};
