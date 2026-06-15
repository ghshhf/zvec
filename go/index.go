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
	IndexTypeIVF    IndexType = C.ZVEC_INDEX_TYPE_IVF
	IndexTypeHNSW   IndexType = C.ZVEC_INDEX_TYPE_HNSW
	IndexTypeFlat   IndexType = C.ZVEC_INDEX_TYPE_FLAT
	IndexTypeFTS    IndexType = C.ZVEC_INDEX_TYPE_FTS
	IndexTypeVamana IndexType = C.ZVEC_INDEX_TYPE_VAMANA
)

// MetricType represents the distance metric.
type MetricType int

const (
	MetricTypeL2     MetricType = C.ZVEC_METRIC_L2
	MetricTypeIP     MetricType = C.ZVEC_METRIC_IP
	MetricTypeCosine MetricType = C.ZVEC_METRIC_COSINE
)

// QueryParams holds query parameters for different index types.
type QueryParams struct {
	HNSW *HNSWQueryParams
	IVF  *IVFQueryParams
	Flat *FlatQueryParams
}

// HNSWQueryParams holds HNSW-specific query parameters.
type HNSWQueryParams struct {
	EFSearch      int
	Radius        float32
	IsLinear      bool
	IsUsingRefiner bool
}

// IVFQueryParams holds IVF-specific query parameters.
type IVFQueryParams struct {
	NProbe         int
	ScaleFactor   float32
	Radius        float32
	IsLinear      bool
	IsUsingRefiner bool
}

// FlatQueryParams holds Flat-specific query parameters.
type FlatQueryParams struct {
	ScaleFactor   float32
	Radius        float32
	IsLinear      bool
	IsUsingRefiner bool
}

// CreateIndex creates a vector index on a field.
func CreateIndex(coll *Collection, fieldName string, indexType IndexType, metric MetricType) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}

	fieldStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldStr))

	indexSchema := C.zvec_collection_schema_create_index(fieldStr, C.zvec_index_type_t(indexType))
	if indexSchema == nil {
		return fmt.Errorf("failed to create index schema")
	}
	defer C.zvec_collection_schema_destroy_index(indexSchema)

	code := C.zvec_collection_create_index(coll.handle, indexSchema)
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to create index: %s", GetLastError())
	}

	return nil
}

// DropIndex drops an index from a collection.
func DropIndex(coll *Collection, indexName string) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}

	indexStr := C.CString(indexName)
	defer C.free(unsafe.Pointer(indexStr))

	code := C.zvec_collection_drop_index(coll.handle, indexStr)
	if code != C.ZVEC_OK {
		return fmt.Errorf("failed to drop index: %s", GetLastError())
	}

	return nil
}

// IndexInfo holds information about an index.
type IndexInfo struct {
	Name   string
	Type   IndexType
	Metric MetricType
}

// ListIndexes lists all indexes in a collection.
func ListIndexes(coll *Collection) ([]IndexInfo, error) {
	// TODO: Implement using zvec_collection_get_stats and parsing index names
	return nil, fmt.Errorf("not implemented yet")
}
