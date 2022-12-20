#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

echo "Building Genesis Generator..."

CAMINO_NODE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
source "$CAMINO_NODE_PATH"/scripts/constants.sh

target="$build_dir/genesis-generator"
go build -o "$target" "$CAMINO_NODE_PATH/tools/genesis_generator/"*.go
