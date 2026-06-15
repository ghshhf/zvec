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
	var coll *C.zvec_collection_t
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))

	// Create default options
	opts := C.zvec_collection_options_create()
	if opts == nil {
		return nil, fmt.Errorf("failed to create options")
	}
	defer C.zvec_collection_options_destroy(opts)

	// Set read-only mode
	if readOnly {
		C.zvec_collection_options_set_read_only(opts, C.bool(true))
	}

	code := C.zvec_collection_open(pathStr, opts, &coll)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to open collection: %s", GetLastError())
	}

	return &Collection{handle: coll, path: path}, nil
}

// Close closes the collection and releases resources.
func (c *Collection) Close() error {
	if c.handle == nil {
		return nil
	}

	code := C.zvec_collection_close(c.handle)
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to close collection: %s", GetLastError())
	}

	c.handle = nil
	return nil
}

// DestroyCollection destroys a collection (deletes all data on disk).
func DestroyCollection(path string) error {
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))

	code := C.zvec_collection_destroy(pathStr)
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to destroy collection: %s", GetLastError())
	}
	return nil
}

// DropCollection drops a collection (alias for DestroyCollection).
func DropCollection(path string) error {
	return DestroyCollection(path)
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

// SetReadOnly sets whether the collection is read-only.
func (o *Options) SetReadOnly(readOnly bool) error {
	code := C.zvec_collection_options_set_read_only(o.handle, C.bool(readOnly))
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to set read_only: %s", GetLastError())
	}
	return nil
}

// SetMaxBufferSize sets the maximum buffer size.
func (o *Options) SetMaxBufferSize(size uint64) error {
	code := C.zvec_collection_options_set_max_buffer_size(o.handle, C.uint64_t(size))
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to set max_buffer_size: %s", GetLastError())
	}
	return nil
}

// GetSchema returns the schema of a collection.
func (c *Collection) GetSchema() (*Schema, error) {
	var schema *C.zvec_collection_schema_t
	code := C.zvec_collection_get_schema(c.handle, &schema)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to get schema: %s", GetLastError())
	}
	return &Schema{handle: schema}, nil
}

// GetOptions returns the options of a collection.
func (c *Collection) GetOptions() (*Options, error) {
	var opts *C.zvec_collection_options_t
	code := C.zvec_collection_get_options(c.handle, &opts)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to get options: %s", GetLastError())
	}
	return &Options{handle: opts}, nil
}

// GetStats returns the statistics of a collection.
func (c *Collection) GetStats() (map[string]interface{}, error) {
	var stats *C.zvec_collection_stats_t
	code := C.zvec_collection_get_stats(c.handle, &stats)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to get stats: %s", GetLastError())
	}

	// TODO: Extract stats into a Go map
	// This requires understanding the full stats structure
	return nil, fmt.Errorf("GetStats not fully implemented yet")
}
