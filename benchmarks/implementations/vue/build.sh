#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/vue/

export NODE_ENV="production"

mkdir -p dist/
pnpm build
