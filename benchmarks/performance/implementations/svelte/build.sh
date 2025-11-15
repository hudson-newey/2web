#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/svelte/

export NODE_ENV="production"

mkdir -p dist/
pnpm build
