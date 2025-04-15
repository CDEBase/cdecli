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

# Build releases for different platforms
PLATFORMS=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
APP_NAME="cli"  # Changed from cde-extension-tool to cli to match existing assets

# Create release directory
mkdir -p ./release

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  
  # Build versioned binary with cli_ prefix
  OUTPUT_NAME=cli_${GOOS}_${GOARCH}
  if [ $GOOS = "windows" ]; then
    OUTPUT_NAME+='.exe'
  fi
  echo "Building versioned for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$VERSION" -o ./release/$OUTPUT_NAME ./src/cli
  
  # Create a copy with version in name
  VERSIONED_NAME=cli_${VERSION}_${GOOS}_${GOARCH}
  if [ $GOOS = "windows" ]; then
    VERSIONED_NAME+='.exe'
  fi
  cp ./release/$OUTPUT_NAME ./release/$VERSIONED_NAME
done

# Create source code archives
git archive --format=zip --output=./release/Source\ code.zip HEAD
git archive --format=tar.gz --output=./release/Source\ code.tar.gz HEAD

echo "Release v$VERSION created!"
