#!/usr/bin/env bash
set -euo pipefail

# We dynamically generate the 2web project when we run the benchmarks so that we
# always have the latest version.

2web new ./2web/project/

source ./.env
web-codegen-scorer eval --env=2web/config.mjs
