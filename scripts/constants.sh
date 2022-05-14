#!/usr/bin/env bash
#
# Use lower_case variables in the scripts and UPPER_CASE variables for override
# Use the constants.sh for env overrides
# Use the versions.sh to specify versions
#

# Set the PATHS
GOPATH="$(go env GOPATH)"

# Where CaminoGo binary goes
build_dir="$CAMINO_NODE_PATH/build"
camino_node_path="$build_dir/camino-node"
plugin_dir="$build_dir/plugins"

# Camino docker hub
# c4tplatform/camino-node - defaults to local as to avoid unintentional pushes
# You should probably set it - export DOCKER_REPO='c4tplatform'
camino_node_dockerhub_repo=${DOCKER_REPO:-"c4tplatform"}"/camino-node"

# Current branch
current_branch=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match || true)

git_commit=${CAMINO_NODE_COMMIT:-$( git rev-list -1 HEAD )}

# caminoethvm git tag and sha
module=$(grep caminoethvm $CAMINO_NODE_PATH/go.mod)
# trim leading
module="${module#"${module%%[![:space:]]*}"}"
t=(${module//\ / })
caminoethvm_tag=${t[1]}

c=$(git ls-remote https://$module);c=(${c//\ / })
caminoethvm_commit=${c[0]}

# Static compilation
static_ld_flags=''
if [ "${STATIC_COMPILATION:-}" = 1 ]
then
    export CC=musl-gcc
    which $CC > /dev/null || ( echo $CC must be available for static compilation && exit 1 )
    static_ld_flags=' -extldflags "-static" -linkmode external '
fi
