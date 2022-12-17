#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go test ${1-} -timeout="120s" -coverprofile="coverage.out" -covermode="atomic" $(go list ./... | grep -v /mocks | grep -v proto | grep -v tests | grep -v /dependencies/)
