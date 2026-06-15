// Package zvec provides Go bindings for the zvec vector database.
// This file implements batch operations and utility functions.
package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// ==================== Batch Operations ====================

// Note: Batch operations are inherently supported by the C API.
// InsertDocuments(), DeleteDocuments(), etc. all accept arrays.
// The following are convenience wrappers for very large datasets
// that need to be processed in smaller batches.

// BatchInsert inserts documents in batches.
// If batchSize <= 0, all documents are inserted in a single batch.
// This is a convenience wrapper around InsertDocuments().
func BatchInsert(coll *Collection, docs []Document, batchSize int) error {
	if len(docs) == 0 {
		return nil
	}

	if batchSize <= 0 || batchSize >= len(docs) {
		// Insert all at once
		_, err := InsertDocuments(coll, docs)
		return err
	}

	// Insert in batches
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}

		_, err := InsertDocuments(coll, docs[i:end])
		if err != nil {
			return fmt.Errorf("batch insert failed at batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}

// BatchSearch performs multiple vector searches in batch.
// This is more efficient than searching one by one when you have multiple queries.
// Note: The C API does not have a native batch search function,
// so this is a convenience wrapper that calls Search() multiple times.
func BatchSearch(coll *Collection, queries [][]float32, fieldName string, topK int,
	params *QueryParams) ([][]SearchResult, error) {

	results := make([][]SearchResult, len(queries))
	for i, query := range queries {
		r, err := Search(coll, fieldName, query, topK, "", nil, params)
		if err != nil {
			return nil, fmt.Errorf("batch search failed at query %d: %w", i, err)
		}
		results[i] = r
	}

	return results, nil
}

// ==================== Collection Management ====================

// FlushCollection flushes data to disk.
// Wrapper around Flush() for API consistency.
func FlushCollection(coll *Collection) error {
	return Flush(coll)
}

// CompactCollection compacts the collection data files.
// Wrapper around Compact() for API consistency.
func CompactCollection(coll *Collection) error {
	return Compact(coll)
}

// GetCollectionInfoFromPath returns information about a collection at the given path.
// Note: This opens the collection, reads stats, then closes it.
func GetCollectionInfoFromPath(path string) (string, error) {
	coll, err := OpenCollection(path, true) // read-only
	if err != nil {
		return "", fmt.Errorf("failed to open collection: %w", err)
	}
	defer coll.Close()

	return GetCollectionInfo(coll)
}

// ==================== Schema Utilities ====================

// ValidateSchema validates a schema.
// Returns nil if valid, or an error describing the issue.
func ValidateSchema(schema *Schema) error {
	if schema == nil || schema.handle == nil {
		return fmt.Errorf("schema is nil")
	}

	code := C.zvec_collection_schema_validate(schema.handle)
	if code != C.ZVEC_OK {
		return fmt.Errorf("schema validation failed: %s", GetLastError())
	}
	return nil
}

// HasField checks if a field exists in the schema.
func HasField(schema *Schema, fieldName string) bool {
	if schema == nil || schema.handle == nil {
		return false
	}

	fieldStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldStr))

	return bool(C.zvec_collection_schema_has_field(schema.handle, fieldStr))
}

// HasIndex checks if an index exists in the schema.
func HasIndex(schema *Schema, indexName string) bool {
	if schema == nil || schema.handle == nil {
		return false
	}

	indexStr := C.CString(indexName)
	defer C.free(unsafe.Pointer(indexStr))

	return bool(C.zvec_collection_schema_has_index(schema.handle, indexStr))
}

// ==================== Global Configuration ====================

// Config holds global zvec configuration.
type Config struct {
	DataConfig *DataConfig
	LogConfig *LogConfig
}

// DataConfig holds data-related configuration.
type DataConfig struct {
	MemoryLimitMB  uint64
	QueryThreadCount int
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level    LogLevel
	FilePath string
}

// Init initializes zvec with the given configuration.
func Init(config *Config) error {
	if config == nil {
		return nil
	}

	// Set query thread count (if specified)
	if config.DataConfig != nil && config.DataConfig.QueryThreadCount > 0 {
		// TODO: Call C API when available
		// C.zvec_config_data_set_query_thread_count(...)
	}

	// Set log level (if specified)
	if config.LogConfig != nil {
		_ = SetLogLevel(config.LogConfig.Level) // Best effort
	}

	return nil
}

// Cleanup cleans up zvec resources.
// Call this when shutting down to release resources.
func Cleanup() {
	C.zvec_clear_error() // Clear any pending errors
	// TODO: Add more cleanup when C API provides it
}
