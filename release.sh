#!/bin/bash
set -ex

# Ensure we have a version argument
if [ -z "$1" ]; then
  echo "Usage: ./release.sh VERSION"
  exit 1
fi

VERSION=$1

# Create version tag
git tag $VERSION -a -m "release v$VERSION"
git tag latest -f -a -m "release v$VERSION"

# Push tags
git push -f --tags

# Build releases for different platforms (optional)
PLATFORMS=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
APP_NAME="cde-extension-tool"

# Create release directory
mkdir -p ./release

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  OUTPUT_NAME=$APP_NAME-$VERSION-$GOOS-$GOARCH
  
  if [ $GOOS = "windows" ]; then
    OUTPUT_NAME+=".exe"
  fi

  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$VERSION" -o ./release/$OUTPUT_NAME ./src/cli
done

echo "Release v$VERSION created!"
