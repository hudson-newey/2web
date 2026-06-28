# 2Web Compiler

## Command Line Options

- `-i <input_path>` (default: index.html)
- `-o <output_path>` (default: dist/index.html)
- `--dev-tools` (default: false) Toggle embedding devtools into the page
- `--production` (default: false) Perform strong optimizations that hurt code readability but also increase speed.
- `--silent` (default: false) Do not write log information
- `--no-cache` (default: false) Disable all 2web build asset caching
- `--stdin` (default: false) Read from STDIN instead from the input_path
- `--stdout` (default: false) Write to STDOUT instead of writing to the output_path
- `--verbose` (default: false) Print extra debug information to the console
- `--no-runtime-optimizations` (default: false) Do not ship additional runtime optimizations
- `--format` (default: false) Formats output assets for readability
- `--ignore-errors` (default: false) Ignores errors in production builds. This allows you to ship compiler errors.
- `--serial` (default: false) Disables parallel compilation of assets
- `--dry-run` (default: false) Simulates a build without writing any files to disk

## Sub Commands

- `cache` (TODO)
  - `cache clean` Deletes all cache files

## Environment Variables

- `__2WEB_CACHE_PATH` (default: ./.cache/) A path that will be used for the build cache

## Production Builds

- Performs code minification
- Performs a handful of runtime [optimizations](../docs/README.md)

## Language Support

❌ = Not working, 🔧 = Developer preview, ✅ = Production ready, ➖ = Not Planned

### Supported Scripting Languages

| Package                                      | State | JS Target | WASM Target | Server Target |
| -------------------------------------------- | ----- | --------- | ----------- | ------------- |
| [Nim](https://nim-lang.org)                  | ❌     | ❌         | ❌           | ❌             |
| [Gleam](https://gleam.run)                   | ❌     | ❌         | ❌           | ❌             |
| [F#](https://fsharp.org)                     | ❌     | ❌         | ❌           | ❌             |
| JavaScript                                   | ✅     | ✅         | ➖           | ❌             |
| [TypeScript](https://www.typescriptlang.org) | ✅     | ✅         | ➖           | ❌             |
| [Flow](https://flow.org)                     | ❌     | ❌         | ➖           | ❌             |
| [CoffeeScript](https://coffeescript.org)     | ❌     | ❌         | ➖           | ❌             |
| [Elm](https://elm-lang.org)                  | ❌     | ❌         | ➖           | ❌             |
| [clojurescript](https://clojurescript.org)   | ❌     | ❌         | ➖           | ❌             |
| [reason](https://reasonml.github.io)         | ❌     | ❌         | ➖           | ❌             |
| [rescript](https://rescript-lang.org)        | ❌     | ❌         | ➖           | ❌             |
| [purescript](https://www.purescript.org)     | ❌     | ❌         | ➖           | ❌             |
| [Civet](https://civet.dev)                   | ❌     | ❌         | ➖           | ❌             |
| [Rust](https://rust-lang.org)                | ❌     | ➖         | ❌           | ❌             |
| C/C++                                        | ❌     | ➖         | ❌           | ❌             |
| C#                                           | ❌     | ➖         | ❌           | ❌             |
| VB.net                                       | ❌     | ➖         | ❌           | ❌             |

### Supported Markup Languages

| Package | State | File Extensions |
| ------- | ----- | --------------- |
| html    | ✅     | `.html`, `.htm` |
| xhtml   | ✅     | `.xhtml`        |
| xml     | ✅     | `.xml`          |
| xslt    | ✅     | `.xslt`, `.xsl` |
| txt     | ✅     | `.txt`          |
| pdf     | ✅     | `.pdf`          |
| md      | ✅     | `.md`           |
| pug     | ❌     | `.pug`          |
| tex     | ❌     | `.tex`          |
| org     | ❌     | `.org`          |
| docx    | 🔧     | `.docx`, `.doc` |
| odt     | ✅     | `.odt`          |
| php     | ❌     | `.php`          |

### Supported Style Languages

| Package | State |
| ------- | ----- |
| css     | ✅     |
| sass    | ❌     |
| scss    | ❌     |
| less    | ❌     |

### Supported Frameworks

The following frameworks can be used in component islands.

| Package | State | File Extensions            |
| ------- | ----- | -------------------------- |
| React   | ❌     | `.react.jsx`, `.react.tsx` |
| Vue     | ❌     | `.vue`                     |
| Svelte  | ❌     | `.svelte`                  |
| hugo    | ❌     | `.hugo.html`               |

## Compiler Developer Options

These flags are not intended to be used by compiler consumers.
They provide some nice debug information that can be helpful when developing the
compiler.

### Developer CLI Options

- `--verbose-lexer` (default: false) Logs the lexer output to the console
- `--verbose-ast` (default: false) Logs the abstract syntax tree to the console
- `--listen` (default: false) Instead of automatically exiting after compilation, open a websocket connection that can accept commands.

### Listener Commands

If the `--listen` flag is used, the compiler does not automatically exit and
will instead open up websocket endpoints where you can send commands.
Using listener commands can improve the performance of the 2web compiler since
the compiler doesn't need to exit / reopen after every compilation.

| Command   | Description             | Action | Body   |
| --------- | ----------------------- | ------ | ------ |
| Recompile | Re-compiles application | 0o1    | _N.A._ |
