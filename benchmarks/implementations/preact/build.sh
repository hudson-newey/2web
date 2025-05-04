#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/preact/

mkdir -p dist/
pnpm build
