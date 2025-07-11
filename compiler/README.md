# 2Web Compiler

## Command Line Options

- `-i <input_path>` (default: index.html)
- `-o <output_path>` (default: dist/index.html)
- `--dev-tools` (default: false) Toggle embedding devtools into the page
- `--production` (default: false) Perform strong optimizations that hurt code readability but also increase speed.
- `--silent` (default: false) Do not write log information
- `--stdin` (default: false) Read from STDIN instead from the input_path
- `--stdout` (default: false) Write to STDOUT instead of writing to the output_path

### Dev Tools

Placeholder (not fully implemented)

### Production Builds

- Performs code minification
- Performs a handful of runtime [optimizations](../docs/README.md)

## More Information

[GitHub](https://github.com/hudson-newey/2web)
