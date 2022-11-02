#!/usr/bin/env bash
#
# Use lower_case variables in the scripts and UPPER_CASE variables for override
# Use the constants.sh for env overrides
# Use the versions.sh to specify versions
#

# Set the PATHS
GOPATH="$(go env GOPATH)"

# Set the CGO flags to use the portable version of BLST
#
# We use "export" here instead of just setting a bash variable because we need
# to pass this flag to all child processes spawned by the shell.
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"

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
module=$(grep -m1 caminoethvm $CAMINO_NODE_PATH/go.mod)
replace=$(echo $module | cut -d' ' -f4)
if [ ${replace} ]
then
    module=$replace
fi

# trim leading
module="${module#"${module%%[![:space:]]*}"}"
t=(${module//\ / })
echo ${#t[*]}
if [ ${#t[*]} -gt 1 ]
then
    caminoethvm_tag=${t[1]}

    c=$(git ls-remote https://$module);c=(${c//\ / })
    caminoethvm_commit=${c[0]}
else
    caminoethvm_tag="v0.0.0"
    caminoethvm_commit="undefined"
fi

# Static compilation
static_ld_flags=''
if [ "${STATIC_COMPILATION:-}" = 1 ]
then
    export CC=musl-gcc
    which $CC > /dev/null || ( echo $CC must be available for static compilation && exit 1 )
    static_ld_flags=' -extldflags "-static" -linkmode external '
fi
