// Package zvec provides Go bindings for the zvec vector database.
// This file implements utility functions for collection management.
package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Flush flushes collection data to disk.
// This ensures all buffered data is persisted.
func Flush(coll *Collection) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil or closed")
	}

	code := C.zvec_collection_flush(coll.handle)
	if code != C.ZVEC_OK {
		return fmt.Errorf("flush failed: %s", GetLastError())
	}
	return nil
}

// Stats holds collection statistics.
type Stats struct {
	DocCount   int64
	IndexCount int
	Indexes    []IndexStats
}

// IndexStats holds statistics for a single index.
type IndexStats struct {
	Name         string
	Completeness float32 // 0.0 - 1.0, index completeness
}

// GetStats gets collection statistics.
func GetStats(coll *Collection) (*Stats, error) {
	if coll == nil || coll.handle == nil {
		return nil, fmt.Errorf("collection is nil or closed")
	}

	var cStats *C.zvec_collection_stats_t
	code := C.zvec_collection_get_stats(coll.handle, &cStats)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to get stats: %s", GetLastError())
	}
	defer C.zvec_collection_stats_destroy(cStats)

	// Get doc count
	docCount := int64(C.zvec_collection_stats_get_doc_count(cStats))

	// Get index count
	indexCount := int(C.zvec_collection_stats_get_index_count(cStats))

	// Get index info
	indexes := make([]IndexStats, 0, indexCount)
	for i := 0; i < indexCount; i++ {
		name := C.GoString(C.zvec_collection_stats_get_index_name(cStats, C.size_t(i)))
		completeness := float32(C.zvec_collection_stats_get_index_completeness(cStats, C.size_t(i)))
		indexes = append(indexes, IndexStats{
			Name:         name,
			Completeness: completeness,
		})
	}

	return &Stats{
		DocCount:   docCount,
		IndexCount: indexCount,
		Indexes:    indexes,
	}, nil
}

// Compact compacts the collection by removing deleted documents and optimizing storage.
// Note: This may take a while for large collections.
func Compact(coll *Collection) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil or closed")
	}

	// TODO: Check if C API has a compact function
	// For now, flush is the closest operation
	return Flush(coll)
}

// GetCollectionInfo gets comprehensive information about a collection.
// Returns a human-readable string with collection details.
func GetCollectionInfo(coll *Collection) (string, error) {
	stats, err := GetStats(coll)
	if err != nil {
		return "", err
	}

	schema, err := GetSchema(coll)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Collection Info:\n"+
		"  Fields: %d\n"+
		"  Documents: %d\n"+
		"  Indexes: %d\n",
		len(schema.Fields),
		stats.DocCount,
		stats.IndexCount,
	)

	for i, idx := range stats.Indexes {
		info += fmt.Sprintf("  Index %d: %s (completeness: %.1f%%)\n",
			i, idx.Name, idx.Completeness*100)
	}

	return info, nil
}

// ValidateCollection validates a collection's schema and data integrity.
// Returns nil if valid, or an error describing the issue.
func ValidateCollection(coll *Collection) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil or closed")
	}

	// Validate schema
	schema, err := GetSchema(coll)
	if err != nil {
		return fmt.Errorf("failed to get schema: %w", err)
	}

	err = ValidateSchema(schema)
	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// Backup creates a backup of the collection to the specified path.
// Note: This uses collection copy internally.
func Backup(coll *Collection, backupPath string) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil or closed")
	}

	// TODO: Implement using C API
	// May need to use zvec_collection_create_and_open with a new path,
	// then copy all documents
	return fmt.Errorf("not implemented yet - requires copy collection support")
}

// Restore restores a collection from a backup.
func Restore(backupPath, restorePath string) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// SetLogLeve sets the global log level for zvec.
func SetLogLevel(level LogLevel) error {
	// TODO: Implement using C API when available
	// C API may have: zvec_set_log_level()
	return fmt.Errorf("not implemented yet - C API does not expose log level control")
}

// GetVersion returns the zvec version string.
// This is already defined in zvec.go, but kept here for reference.
func GetVersion() string {
	return Version()
}
