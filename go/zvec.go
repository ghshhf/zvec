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

// CheckVersion checks if the runtime version is compatible with the required version.
func CheckVersion(major, minor, patch int) bool {
	return bool(C.zvec_check_version(C.int(major), C.int(minor), C.int(patch)))
}

// GetLastError returns the last error message.
func GetLastError() string {
	var errMsg *C.char
	C.zvec_get_last_error(&errMsg)
	if errMsg == nil {
		return ""
	}
	return C.GoString(errMsg)
}

// ClearError clears the last error.
func ClearError() {
	C.zvec_clear_error()
}

// LogLevel represents the logging level.
type LogLevel int

const (
	LogLevelDebug LogLevel = C.ZVEC_LOG_DEBUG
	LogLevelInfo  LogLevel = C.ZVEC_LOG_INFO
	LogLevelWarn  LogLevel = C.ZVEC_LOG_WARN
	LogLevelError LogLevel = C.ZVEC_LOG_ERROR
	LogLevelFatal LogLevel = C.ZVEC_LOG_FATAL
)

// SetLogLevel sets the global log level.
func SetLogLevel(level LogLevel) {
	C.zvec_set_log_level(C.zvec_log_level_t(level))
}

// Collection represents a zvec collection.
type Collection struct {
	handle *C.zvec_collection_t
	path   string
}

// Schema represents a collection schema.
type Schema struct {
	handle *C.zvec_collection_schema_t
}

// Options represents collection options.
type Options struct {
	handle *C.zvec_collection_options_t
}

// CreateCollection creates a new collection with the given schema and options.
func CreateCollection(path string, schema *Schema, opts *Options) (*Collection, error) {
	var coll *C.zvec_collection_t
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))

	var cOpts *C.zvec_collection_options_t
	if opts != nil {
		cOpts = opts.handle
	}

	code := C.zvec_collection_create_and_open(pathStr, schema.handle, cOpts, &coll)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to create collection: %s", GetLastError())
	}

	return &Collection{handle: coll, path: path}, nil
}

// OpenCollection opens an existing collection.
func OpenCollection(path string, readOnly bool) (*Collection, error) {
	// TODO: Implement using C API - need to check the exact C function name
	return nil, fmt.Errorf("not implemented yet")
}

// Close closes the collection.
func (c *Collection) Close() error {
	if c.handle == nil {
		return nil
	}
	// TODO: Implement using C API
	c.handle = nil
	return nil
}

// CreateSchema creates a new collection schema with the given name.
func CreateSchema(name string) (*Schema, error) {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	schema := C.zvec_collection_schema_create(nameStr)
	if schema == nil {
		return nil, fmt.Errorf("failed to create schema")
	}

	return &Schema{handle: schema}, nil
}

// SetName sets the name of the schema.
func (s *Schema) SetName(name string) error {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	code := C.zvec_collection_schema_set_name(s.handle, nameStr)
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to set schema name: %s", GetLastError())
	}
	return nil
}

// AddField adds a field to the schema.
func (s *Schema) AddField(name string, fieldType int) error {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	code := C.zvec_collection_schema_add_field(s.handle, nameStr, C.zvec_field_type_t(fieldType))
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to add field: %s", GetLastError())
	}
	return nil
}

// CreateOptions creates new collection options.
func CreateOptions() (*Options, error) {
	opts := C.zvec_collection_options_create()
	if opts == nil {
		return nil, fmt.Errorf("failed to create options")
	}
	return &Options{handle: opts}, nil
}

// SetEnableMMap sets whether to enable mmap.
func (o *Options) SetEnableMMap(enable bool) error {
	code := C.zvec_collection_options_set_enable_mmap(o.handle, C.bool(enable))
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to set enable_mmap: %s", GetLastError())
	}
	return nil
}
