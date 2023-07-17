#!/bin/sh
export GOOS=darwin
ARCH=$(arch)
# Apple no longer support 32bit apps, so as Go 1.15+
if [ "$ARCH" = "i386" ]; then
  ARCH="amd64"
fi
export GOARCH="$ARCH"

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse HEAD)

TARGET="build/onf-${GOOS}-${GOARCH}-v${VERSION}"
FLAGS="-X main.version=${VERSION} -X main.gitCommit=${COMMIT}"

echo "Building $TARGET"
go build -o "${TARGET}" -ldflags "${FLAGS}" ./cmd
