#!/usr/bin/env bash
# Build script for Go coreutils

set -e

echo "Building coreutils..."
go build -ldflags="-s -w" -o coreutils .
echo "Built: ./coreutils ($(du -h coreutils | cut -f1))"

# Optional: create symlinks
if [ "$1" = "--symlinks" ]; then
    BINDIR="${2:-./bin}"
    mkdir -p "$BINDIR"
    # Get list of commands
    CMDS=$(./coreutils 2>&1 | grep "^  " | awk '{print $1}')
    for cmd in $CMDS; do
        ln -sf "$(realpath coreutils)" "$BINDIR/$cmd"
        echo "  -> $BINDIR/$cmd"
    done
    echo "Symlinks created in $BINDIR/"
fi
