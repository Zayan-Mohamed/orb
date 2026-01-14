#!/bin/bash

# Orb Cross-Platform Build Script
# Builds static binaries for Linux, macOS, and Windows

set -e

VERSION=${VERSION:-"dev"}
BUILD_DIR="build"
LDFLAGS="-s -w -X main.version=${VERSION}"

echo "Building Orb v${VERSION}..."
echo ""

# Clean build directory
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build for each platform
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    output_name="orb-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="${LDFLAGS}" \
        -o ${BUILD_DIR}/${output_name} \
        .
    
    echo "  ✓ ${BUILD_DIR}/${output_name}"
done

echo ""
echo "Build complete!"
echo ""
echo "Binaries in ${BUILD_DIR}/"
ls -lh ${BUILD_DIR}/
