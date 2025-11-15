#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/react/

export NODE_ENV="production"

mkdir -p dist/
pnpm build
