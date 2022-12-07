#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Camino-Node root folder
CAMINO_NODE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$CAMINO_NODE_PATH"/scripts/constants.sh

# Download dependencies
echo "Initializing git submodules..."
git --git-dir $CAMINO_NODE_PATH/.git submodule update --init

echo "Downloading dependencies..."
(cd $CAMINO_NODE_PATH && go mod download)

# Build caminogo
"$CAMINO_NODE_PATH"/scripts/build_camino.sh

# Make plugin folder
mkdir -p $plugin_dir

# Exit build successfully if the binaries are created
if [[ -f "$camino_node_path" ]]; then
        echo "Build Successful"
        exit 0
else
        echo "Build failure" >&2
        exit 1
fi
