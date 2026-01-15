#!/bin/bash

# Orb Cross-Platform Build Script
# Builds static binaries for Linux, macOS, and Windows

set -e

VERSION=${VERSION:-"dev"}
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_DIR="build"
LDFLAGS="-s -w -X github.com/Zayan-Mohamed/orb/cmd.Version=${VERSION} -X github.com/Zayan-Mohamed/orb/cmd.GitCommit=${GIT_COMMIT} -X github.com/Zayan-Mohamed/orb/cmd.BuildDate=${BUILD_DATE}"

echo "Building Orb v${VERSION}..."
echo "Git Commit: ${GIT_COMMIT}"
echo "Build Date: ${BUILD_DATE}"
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
