package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// IndexType represents the type of index.
type IndexType int

const (
	IndexTypeIVF   IndexType = C.ZVEC_INDEX_TYPE_IVF
	IndexTypeHNSW  IndexType = C.ZVEC_INDEX_TYPE_HNSW
	IndexTypeFlat  IndexType = C.ZVEC_INDEX_TYPE_FLAT
	IndexTypeFTS   IndexType = C.ZVEC_INDEX_TYPE_FTS
	IndexTypeVamana IndexType = C.ZVEC_INDEX_TYPE_VAMANA
)

// MetricType represents the distance metric.
type MetricType int

const (
	MetricTypeL2     MetricType = C.ZVEC_METRIC_L2
	MetricTypeIP     MetricType = C.ZVEC_METRIC_IP
	MetricTypeCosine MetricType = C.ZVEC_METRIC_COSINE
)

// CreateIndex creates a vector index on a field.
func CreateIndex(coll *Collection, fieldName string, indexType IndexType, metric MetricType) error {
	fieldStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldStr))

	// TODO: Need to check the exact C API function name and signature
	// code := C.zvec_collection_create_index(coll.handle, fieldStr, ...)
	return fmt.Errorf("not implemented yet")
}

// DropIndex drops an index from a collection.
func DropIndex(coll *Collection, indexName string) error {
	indexStr := C.CString(indexName)
	defer C.free(unsafe.Pointer(indexStr))

	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// IndexInfo holds information about an index.
type IndexInfo struct {
	Name   string
	Type   IndexType
	Metric MetricType
}
