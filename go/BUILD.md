# Go SDK Build and Test Guide

This document explains how to build the zvec C library and run Go SDK integration tests.

## Prerequisites

### On Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install -y build-essential cmake
```

### On macOS
```bash
brew install cmake gcc
```

### On Windows
1. Install [MSYS2](https://www.msys2.org/)
2. Open MSYS2 MinGW 64-bit terminal
3. Run:
```bash
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake make
```

## Step 1: Build zvec C Library

```bash
cd /path/to/zvec
mkdir -p build
cd build
cmake ..
make -j$(nproc)
```

This will produce:
- Linux: `libzvec.so`in `build/`
- macOS: `libzvec.dylib` in `build/`
- Windows: `zvec.dll` in `build/`

## Step 2: Copy Library to Go SDK

```bash
# Linux
cp build/libzvec.so go/lib/

# macOS
cp build/libzvec.dylib go/lib/

# Windows
cp build/zvec.dll go/lib/
```

Also copy the header file:
```bash
cp src/include/zvec/c_api.h go/include/zvec/
```

## Step 3: Update Go bindings (if needed)

The Go SDK uses cgo to link against the C library. Update `go/zvec.go`:

```go
// Update the cgo LDFLAGS to point to the correct library path
/*
#cgo linux LDFLAGS: -L${SRCDIR}/lib -lzvec -lm
#cgo darwin LDFLAGS: -L${SRCDIR}/lib -lzvec -lm
#cgo windows LDFLAGS: -L${SRCDIR}/lib -lzvec -lm
#include <zvec/c_api.h>
*/
import "C"
```

## Step 4: Run Integration Tests

```bash
cd go
go test -v -run TestIntegration
```

## Troubleshooting

### "cgo: C compiler not found"
Install a C compiler:
- Linux: `sudo apt-get install gcc`
- macOS: `xcode-select --install`
- Windows: Install MinGW-w64 via MSYS2

### "undefined reference to zvec_*"
The C library is not linked correctly. Check:
1. Library path in `#cgo LDFLAGS`
2. Library file exists in `go/lib/`
3. Header file exists in `go/include/zvec/`

### "dyld: Library not loaded" (macOS)
Set the library path:
```bash
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$(pwd)/go/lib
```

### "cannot open shared object file" (Linux)
Set the library path:
```bash
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(pwd)/go/lib
```
