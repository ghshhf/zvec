// Package zvec provides Go bindings for the zvec vector database.
// This file implements vector search functionality.
package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// SearchResult represents a single search result.
type SearchResult struct {
	Doc      Document  // The document
	Score    float32   // Similarity score
	Distance float32   // Vector distance
}

// Search performs vector similarity search.
// fieldName: name of the vector field to search
// queryVector: query vector (float32 array)
// topK: number of results to return
// filter: optional filter expression (can be "")
// outputFields: fields to include in results (nil = all fields)
// params: query parameters for specific index type (can be nil for defaults)
func Search(coll *Collection, fieldName string, queryVector []float32, topK int,
	filter string, outputFields []string, params *QueryParams) ([]SearchResult, error) {

	if coll == nil || coll.handle == nil {
		return nil, fmt.Errorf("collection is nil or closed")
	}

	if len(queryVector) == 0 {
		return nil, fmt.Errorf("query vector is empty")
	}

	if topK <= 0 {
		return nil, fmt.Errorf("topK must be positive")
	}

	// 1. Create vector query
	query := C.zvec_vector_query_create()
	if query == nil {
		return nil, fmt.Errorf("failed to create vector query: %s", GetLastError())
	}
	defer C.zvec_vector_query_destroy(query)

	// 2. Set field name
	cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(cFieldName))
	code := C.zvec_vector_query_set_field_name(query, cFieldName)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to set field name: %s", GetLastError())
	}

	// 3. Set query vector (float32)
	// size = len(queryVector) * 4 (sizeof(float32))
	vectorSize := C.size_t(len(queryVector) * 4)
	code = C.zvec_vector_query_set_query_vector(
		query,
		unsafe.Pointer(&queryVector[0]),
		vectorSize,
	)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("failed to set query vector: %s", GetLastError())
	}

	// 4. Set topK
	C.zvec_vector_query_set_topk(query, C.int(topK))

	// 5. Set filter (if provided)
	if filter != "" {
		cFilter := C.CString(filter)
		defer C.free(unsafe.Pointer(cFilter))
		code = C.zvec_vector_query_set_filter(query, cFilter)
		if code != C.ZVEC_OK {
			return nil, fmt.Errorf("failed to set filter: %s", GetLastError())
		}
	}

	// 6. Set output fields (if provided)
	if len(outputFields) > 0 {
		cFields := make([]*C.char, len(outputFields)+1) // +1 for nil terminator
		for i, f := range outputFields {
			cFields[i] = C.CString(f)
			defer C.free(unsafe.Pointer(cFields[i]))
		}
		cFields[len(outputFields)] = nil

		code = C.zvec_vector_query_set_output_fields(
			query,
			(**C.char)(unsafe.Pointer(&cFields[0])),
			C.size_t(len(outputFields)),
		)
		if code != C.ZVEC_OK {
			return nil, fmt.Errorf("failed to set output fields: %s", GetLastError())
		}
	}

	// 7. Set index-specific query params (if provided)
	if params != nil {
		if err := setQueryParams(query, params); err != nil {
			return nil, err
		}
	}

	// 8. Execute query
	var cResults **C.zvec_doc_t
	var resultCount C.size_t

	code = C.zvec_collection_query(
		coll.handle,
		query,
		&cResults,
		&resultCount,
	)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("search failed: %s", GetLastError())
	}
	defer C.zvec_docs_free(cResults, resultCount)

	// 9. Convert results to Go structs
	results := make([]SearchResult, 0, int(resultCount))
	for i := 0; i < int(resultCount); i++ {
		cDoc := *(**C.zvec_doc_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cResults)) + uintptr(i)*unsafe.Sizeof(*cResults),
		))

		doc, err := cDocToGo(cDoc)
		if err != nil {
			continue // Skip invalid documents
		}

		// TODO: Get score and distance from result
		// (Need to check C API for score/distance access)
		result := SearchResult{
			Doc:      *doc,
			Score:    0.0, // TODO: Get from C result
			Distance: 0.0, // TODO: Get from C result
		}
		results = append(results, result)
	}

	return results, nil
}

// setQueryParams sets index-specific query parameters.
func setQueryParams(query *C.zvec_vector_query_t, params *QueryParams) error {
	if params == nil {
		return nil
	}

	// HNSW params
	if params.HNSW != nil {
		cParams := C.zvec_query_params_hnsw_create(
			C.bool(params.HNSW.IsUsingRefiner),
		)
		if cParams == nil {
			return fmt.Errorf("failed to create HNSW query params")
		}
		defer C.zvec_query_params_hnsw_destroy(cParams)

		if params.HNSW.EFSearch > 0 {
			C.zvec_query_params_hnsw_set_ef(cParams, C.int(params.HNSW.EFSearch))
		}
		if params.HNSW.Radius > 0 {
			C.zvec_query_params_hnsw_set_radius(cParams, C.float(params.HNSW.Radius))
		}
		C.zvec_query_params_hnsw_set_is_linear(cParams, C.bool(params.HNSW.IsLinear))

		code := C.zvec_vector_query_set_hnsw_params(query, cParams)
		if code != C.ZVEC_OK {
			return fmt.Errorf("failed to set HNSW params: %s", GetLastError())
		}
	}

	// IVF params
	if params.IVF != nil {
		cParams := C.zvec_query_params_ivf_create(
			C.bool(params.IVF.IsUsingRefiner),
		)
		if cParams == nil {
			return fmt.Errorf("failed to create IVF query params")
		}
		defer C.zvec_query_params_ivf_destroy(cParams)

		if params.IVF.NProbe > 0 {
			C.zvec_query_params_ivf_set_nprobe(cParams, C.int(params.IVF.NProbe))
		}
		if params.IVF.ScaleFactor > 0 {
			C.zvec_query_params_ivf_set_scale_factor(cParams, C.float(params.IVF.ScaleFactor))
		}
		if params.IVF.Radius > 0 {
			C.zvec_query_params_ivf_set_radius(cParams, C.float(params.IVF.Radius))
		}
		C.zvec_query_params_ivf_set_is_linear(cParams, C.bool(params.IVF.IsLinear))

		code := C.zvec_vector_query_set_ivf_params(query, cParams)
		if code != C.ZVEC_OK {
			return fmt.Errorf("failed to set IVF params: %s", GetLastError())
		}
	}

	// Flat params
	if params.Flat != nil {
		cParams := C.zvec_query_params_flat_create(
			C.bool(params.Flat.IsUsingRefiner),
			C.float(params.Flat.ScaleFactor),
		)
		if cParams == nil {
			return fmt.Errorf("failed to create Flat query params")
		}
		defer C.zvec_query_params_flat_destroy(cParams)

		if params.Flat.Radius > 0 {
			C.zvec_query_params_flat_set_radius(cParams, C.float(params.Flat.Radius))
		}
		C.zvec_query_params_flat_set_is_linear(cParams, C.bool(params.Flat.IsLinear))

		code := C.zvec_vector_query_set_flat_params(query, cParams)
		if code != C.ZVEC_OK {
			return fmt.Errorf("failed to set Flat params: %s", GetLastError())
		}
	}

	return nil
}
