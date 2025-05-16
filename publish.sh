#!/bin/bash
set -euo pipefail

# Check if version argument is provided
if [ -z "$1" ]; then
    echo "Error: Please provide a version argument (e.g., patch, minor, major)"
    exit 1
fi

# build the assets that we are going to publish
cd ./cli
task build
cd ..

cd ./compiler
task build
cd ..

cd ./packages
pnpm build
cd ..

# Array of directories containing npm packages
directories=(
    "cli"
    "compiler"
    "packages"
    "sdk"
)

# Iterate over directories
for dir in "${directories[@]}"; do
    # Check if directory exists
    if [ -d "$dir" ]; then
        echo "Processing $dir..."

        # Change to directory
        cd "$dir" || {
            echo "Failed to cd into $dir"
            continue
        }

        # Run npm version with provided argument
        npm version "$1" || {
            echo "Failed to update version in $dir"
            cd ..
            continue
        }

        # Publish to npm
        npm publish || {
            echo "Failed to publish $dir"
            cd ..
            continue
        }

        # Return to original directory
        cd ..
        echo "Successfully processed $dir"
    else
        echo "Directory $dir does not exist"
    fi
done

echo "Completed processing all directories"

echo "Don't forget to publish the two-web/kit to JSR"
