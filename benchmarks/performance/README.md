# Benchmarks

Here I compare 2web against other popular web frameworks on a variety of
metrics.

**These are by no means scientific or accurate benchmarks**, and ONLY serve as a
simple smoke test for my web framework.

## How to run

The benchmark will iterate over all iterations, running their `build.sh`
scripts which must output to a local `dist/` directory under the implementation.

The implementations `dist/` directory will be used for benchmarking.

To run the benchmark script, simply run the following command with
[Deno](https://deno.com/) installed.

```sh
$ pnpm start
>
```

Benchmarking software is written in JavaScript so that it can easily interop
with real browsers to test the speed of each framework.
