#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ -z "${VERSION:-}" ]; then
  VERSION="$(git describe --always --tags)"
fi

cat << EOF
VERSION ${VERSION}
EOF
