#!/usr/bin/env bash
set -euo pipefail

cd ./implementations/2web/

mkdir -p dist/
../../../compiler/bin/2webc -i counter.html -o dist/counter.html --production

# TODO: Implement release build implementation

# # by setting the IS_MONO=1 environment variable, we are able to run these
# # benchmarks during development and without creating a separate 2web release
# if [ $2WEB_IS_MONO -eq 1 ]; then
#   echo "Running development build of 2web..."
#   ../../../compiler/build/2web -i counter.html -o dist/counter.html --production
# else
#   # if the IS_MONO environment variable is not set, we assume that the user who
#   # is running these benchmarks is not wanting to use the monorepo version, and
#   # should instead use the latest release installed on the system.
#   2web -i ./counter.html -o dist/counter.html --production
# fi
