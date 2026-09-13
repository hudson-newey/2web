#!/usr/bin/env bash
# Compiler performance benchmark.
#
# Measures the wall clock time of compiling the dev fixture site and a large
# synthetic site (360 pages). Reports the median of N runs so that scheduler
# noise doesn't drown the signal.
#
# Usage: ./bench.sh [runs]

set -u

COMPILER=${COMPILER:-./bin/2webc}
RUNS=${1:-7}
LARGE_SITE=${LARGE_SITE:-/tmp/bench-perf/site}

median() {
	sort -n | awk '{a[NR]=$1} END {if (NR % 2) print a[(NR+1)/2]; else print (a[NR/2]+a[NR/2+1])/2}'
}

bench() {
	local label=$1 target=$2 extra=$3
	local times=()
	for _ in $(seq 1 "$RUNS"); do
		local out
		out=$(mktemp -d)
		local start end
		start=$(date +%s%N)
		"$COMPILER" -i "$target" -o "$out/" $extra --silent > /dev/null 2>&1
		end=$(date +%s%N)
		rm -rf "$out"
		times+=( $(( (end - start) / 1000000 )) )
	done
	local med
	med=$(printf '%s\n' "${times[@]}" | median)
	local min
	min=$(printf '%s\n' "${times[@]}" | sort -n | head -1)
	printf '%-42s median %6s ms   min %6s ms\n' "$label" "$med" "$min"
}

echo "compiler: $COMPILER (runs: $RUNS)"
echo "---"
bench "dev fixtures (33 pages)" dev "--no-cache"
bench "dev fixtures --serial" dev "--no-cache --serial"
bench "large site (360 pages)" "$LARGE_SITE" "--no-cache"
bench "large site --serial" "$LARGE_SITE" "--no-cache --serial"

# Cached incremental rebuild: run once to populate the cache, then measure.
warm_out=$(mktemp -d)
"$COMPILER" -i "$LARGE_SITE" -o "$warm_out/" --silent > /dev/null 2>&1
bench "large site (360 pages, cached)" "$LARGE_SITE" ""
rm -rf "$warm_out"
