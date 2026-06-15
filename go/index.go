package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Index represents a vector index.
type Index struct {
	handle *C.zvec_index_t
	name   string
}

// CreateIndex creates a vector index on a field.
func CreateIndex(collection *Collection, fieldName string, indexType IndexType, metric MetricType) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// DropIndex drops an index.
func DropIndex(collection *Collection, indexName string) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// IndexType represents the type of index.
type IndexType int

const (
	IVF IndexType = iota
	HNSW
	FLAT
)

// MetricType represents the distance metric.
type MetricType int

const (
	L2 MetricType = iota
	IP
	COSINE
)
