#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/vue

mkdir -p dist/
pnpm build
