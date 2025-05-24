# Benchmarks

Here I compare 2web against other popular web frameworks on a variety of
metrics.

These are by no means scientific tests, but serve as a simple smoke test for my
web framework.

## How to run

The benchmark will iterate over all iterations, running their `build.sh`
scripts which must output to a local `dist/` directory under the implementation.

The implementations `dist/` directory will be used for benchmarking.

To run the benchmark script, simply run the following command with
[Deno](https://deno.com/) installed.

```sh
$ ./bench.ts
>
```

Benchmarking software is written in JavaScript so that it can easily interop
with real browsers to test the speed of each framework.

## Currently Benchmarked Frameworks

| Framework         | Build Size (KB) | Build Time (MS) \* |
| ----------------- | --------------- | ------------------ |
| Vanilla HTML + JS | 0.293           | 0                  |
| 2Web              | 0.473           | 9                  |
| Svelte            | 20.122          | 1480               |
| Preact            | 25.064          | 1577               |
| Vue               | 58.705          | 2988               |
| React             | 186.673         | 1733               |

\* Build times are highly experimental and have high std-error.
