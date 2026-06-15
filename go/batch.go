package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// BatchInsert inserts multiple documents in batch.
// This is more efficient than inserting one by one.
func BatchInsert(coll *Collection, docs []map[string]interface{}, batchSize int) error {
	// TODO: Implement batch insert using C API
	// This should batch documents and call insert multiple times
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// BatchSearch performs multiple searches in batch.
// This is more efficient than searching one by one.
func BatchSearch(coll *Collection, queries [][]float32, fieldName string, limit int) ([][]SearchResult, error) {
	// TODO: Implement batch search using C API
	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// Config holds global zvec configuration.
type Config struct {
	DataConfig   *DataConfig
	LogConfig    *LogConfig
}

// DataConfig holds data-related configuration.
type DataConfig struct {
	MemoryLimitMB uint64
	QueryThreadCount int
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level    LogLevel
	FilePath string
}

// Init initializes zvec with the given configuration.
func Init(config *Config) error {
	// Set log level
	if config != nil && config.LogConfig != nil {
		SetLogLevel(config.LogConfig.Level)
	}

	// TODO: Apply other configuration options using C API
	return nil
}

// Cleanup cleans up zvec resources.
func Cleanup() {
	// TODO: Implement cleanup using C API
	C.zvec_clear_error() // At minimum, clear any pending errors
}

// CompactCollection compacts the collection data files.
func CompactCollection(coll *Collection) error {
	// TODO: Implement using C API - zvec_collection_compact
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// FlushCollection flushes data to disk.
func FlushCollection(coll *Collection) error {
	// TODO: Implement using C API - zvec_collection_flush
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// GetCollectionInfo returns information about a collection without opening it.
func GetCollectionInfo(path string) (map[string]interface{}, error) {
	// TODO: Implement using C API
	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// ValidateSchema validates a schema.
func ValidateSchema(schema *Schema) error {
	code := C.zvec_collection_schema_validate(schema.handle)
	if code != C.ZVEC_OK {
		return fmt.Errorf("schema validation failed: %s", GetLastError())
	}
	return nil
}

// HasField checks if a field exists in the schema.
func HasField(schema *Schema, fieldName string) bool {
	fieldStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldStr))

	return bool(C.zvec_collection_schema_has_field(schema.handle, fieldStr))
}

// HasIndex checks if an index exists in the schema.
func HasIndex(schema *Schema, indexName string) bool {
	indexStr := C.CString(indexName)
	defer C.free(unsafe.Pointer(indexStr))

	return bool(C.zvec_collection_schema_has_index(schema.handle, indexStr))
}
