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

// CreateIndex creates a vector index on a field.
// This is a Go wrapper around the C API zvec_collection_create_index().
func CreateIndex(coll *Collection, fieldName string, indexType IndexType, metric MetricType) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}

	fieldStr := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldStr))

	// Create index schema
	indexSchema := C.zvec_collection_schema_create_index(fieldStr, C.zvec_index_type_t(indexType))
	if indexSchema == nil {
		return fmt.Errorf("failed to create index schema")
	}
	defer C.zvec_collection_schema_destroy_index(indexSchema)

	// Set metric type (if applicable)
	// TODO: Set metric type on index schema based on C API

	// Create the index
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

// QueryParams holds query parameters for vector search.
type QueryParams struct {
	IndexType     IndexType
	HNSWParams   *HNSWQueryParams
	IVFParams     *IVFQueryParams
	FlatParams    *FlatQueryParams
	VamanaParams *VamanaQueryParams
}

// HNSWQueryParams holds HNSW-specific query parameters.
type HNSWQueryParams struct {
	EF      int
	Radius  float32
	IsLinear bool
}

// IVFQueryParams holds IVF-specific query parameters.
type IVFQueryParams struct {
	NProbe      int
	ScaleFactor  float32
	IsLinear     bool
}

// FlatQueryParams holds Flat-specific query parameters.
type FlatQueryParams struct {
	ScaleFactor float32
	Radius      float32
	IsLinear    bool
}

// Search performs vector search in a collection.
// This is a Go wrapper around the C API zvec_collection_query().
func Search(coll *Collection, query []float32, fieldName string, limit int, params *QueryParams) ([]SearchResult, error) {
	if coll == nil || coll.handle == nil {
		return nil, fmt.Errorf("collection is nil")
	}

	// TODO: Implement full search using C API
	// 1. Create query params based on index type
	// 2. Call zvec_collection_query()
	// 3. Parse results into []SearchResult
	// 4. Free C resources

	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// SearchByVector performs vector search with a vector.
func SearchByVector(coll *Collection, vector []float32, fieldName string, limit int) ([]SearchResult, error) {
	return Search(coll, vector, fieldName, limit, nil)
}

// SearchByID performs vector search using an existing document's vector.
func SearchByID(coll *Collection, docID string, fieldName string, limit int) ([]SearchResult, error) {
	// TODO: Implement using C API - fetch doc vector, then search
	return nil, fmt.Errorf("not implemented yet")
}
