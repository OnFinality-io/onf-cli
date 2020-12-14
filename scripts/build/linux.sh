export GOOS=linux
export GOARCH=amd64

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse HEAD)

TARGET="build/onf-${GOOS}-${GOARCH}-v${VERSION}"
FLAGS="-X main.version=${VERSION} -X main.gitCommit=${COMMIT}"

echo "Building $TARGET"
go build -o "${TARGET}" -ldflags "${FLAGS}" ./cmd