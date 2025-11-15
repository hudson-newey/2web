#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/preact/

export NODE_ENV="production"

mkdir -p dist/
pnpm build
