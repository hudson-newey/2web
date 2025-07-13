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

## Environment Variables

- `__2_CACHE_PATH` (default: ./.cache/) A path that will be used for the build cache

### Dev Tools

Placeholder (not fully implemented)

### Production Builds

- Performs code minification
- Performs a handful of runtime [optimizations](../docs/README.md)

## Compiler Development Options

These flags are not intended to be used by compiler consumers.
They provide some nice debug information that can be helpful when developing the
compiler.

- `--verbose-lexer` (default: false) Logs the lexer output to the console
- `--verbose-ast` (default: false) Logs the abstract syntax tree to the console

## More Information

[GitHub](https://github.com/hudson-newey/2web)
