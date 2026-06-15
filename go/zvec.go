// Package zvec provides Go bindings for the zvec vector database.
// It uses cgo to call the zvec C API.
package zvec

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lzvec -L./lib
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Version returns the zvec version string.
func Version() string {
	return C.GoString(C.zvec_version())
}

// NewDB creates a new zvec database instance.
func NewDB(config *Config) (*DB, error) {
	// TODO: Implement based on C API
	return nil, fmt.Errorf("not implemented yet")
}

// Config holds database configuration.
type Config struct {
	DBPath   string
	MaxMemGB int
}

// DB represents a zvec database instance.
type DB struct {
	handle unsafe.Pointer
}
