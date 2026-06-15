#!/bin/bash
# Build script for zvec C library and Go SDK

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ZVEC_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== Building zvec C library ==="
echo "ZVEC_ROOT: $ZVEC_ROOT"

# Create build directory
BUILD_DIR="$ZVEC_ROOT/build"
mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"

# Run cmake
echo "Running cmake..."
cmake ..

# Build
echo "Building..."
make -j$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

# Copy library to Go SDK
echo "=== Copying library to Go SDK ==="
GO_LIB_DIR="$ZVEC_ROOT/go/lib"
mkdir -p "$GO_LIB_DIR"

if [ -f "$BUILD_DIR/libzvec.so" ]; then
    cp "$BUILD_DIR/libzvec.so" "$GO_LIB_DIR/"
    echo "Copied libzvec.so to $GO_LIB_DIR/"
elif [ -f "$BUILD_DIR/libzvec.dylib" ]; then
    cp "$BUILD_DIR/libzvec.dylib" "$GO_LIB_DIR/"
    echo "Copied libzvec.dylib to $GO_LIB_DIR/"
elif [ -f "$BUILD_DIR/zvec.dll" ]; then
    cp "$BUILD_DIR/zvec.dll" "$GO_LIB_DIR/"
    echo "Copied zvec.dll to $GO_LIB_DIR/"
else
    echo "ERROR: Could not find built library in $BUILD_DIR"
    echo "Contents of $BUILD_DIR:"
    ls -la "$BUILD_DIR/"
    exit 1
fi

# Copy header file
echo "=== Copying header files ==="
GO_INCLUDE_DIR="$ZVEC_ROOT/go/include/zvec"
mkdir -p "$GO_INCLUDE_DIR"
cp "$ZVEC_ROOT/src/include/zvec/c_api.h" "$GO_INCLUDE_DIR/"
echo "Copied c_api.h to $GO_INCLUDE_DIR/"

echo ""
echo "=== Build complete! ==="
echo "Library: $GO_LIB_DIR/"
echo "Headers: $GO_INCLUDE_DIR/"
echo ""
echo "Next steps:"
echo "  cd go"
echo "  go test -v"
