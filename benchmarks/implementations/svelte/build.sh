#!/usr/bin/env bash
set -euo pipefail

export NODE_ENV="production"

mkdir -p dist/
pnpm build
