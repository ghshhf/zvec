package zvec

/*
#include <zvec/c_api.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// FieldType represents the type of a field.
type FieldType int

const (
	FieldTypeInt32        FieldType = C.ZVEC_FIELD_TYPE_INT32
	FieldTypeInt64        FieldType = C.ZVEC_FIELD_TYPE_INT64
	FieldTypeFloat        FieldType = C.ZVEC_FIELD_TYPE_FLOAT
	FieldTypeDouble       FieldType = C.ZVEC_FIELD_TYPE_DOUBLE
	FieldTypeString       FieldType = C.ZVEC_FIELD_TYPE_STRING
	FieldTypeBinary       FieldType = C.ZVEC_FIELD_TYPE_BINARY
	FieldTypeBool         FieldType = C.ZVEC_FIELD_TYPE_BOOL
	FieldTypeJSON         FieldType = C.ZVEC_FIELD_TYPE_JSON
	FieldTypeVector       FieldType = C.ZVEC_FIELD_TYPE_VECTOR
	FieldTypeSparseVector FieldType = C.ZVEC_FIELD_TYPE_SPARSE_VECTOR
	FieldTypeDatetime     FieldType = C.ZVEC_FIELD_TYPE_DATETIME
)

// Document represents a document with fields.
type Document struct {
	ID     string
	Fields map[string]interface{}
}

// InsertDocuments inserts documents into a collection.
// This is a Go wrapper around the C API zvec_collection_insert().
func InsertDocuments(coll *Collection, docs []Document) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}
	if len(docs) == 0 {
		return nil
	}

	// Create C documents
	cDocs := make([]*C.zvec_doc_t, len(docs))
	defer func() {
		for _, doc := range cDocs {
			if doc != nil {
				C.zvec_doc_destroy(doc)
			}
		}
	}()

	for i, doc := range docs {
		// Create C document
		cDoc := C.zvec_doc_create()
		if cDoc == nil {
			return fmt.Errorf("failed to create C document for doc %d", i)
		}
		cDocs[i] = cDoc

		// Set PK (primary key = ID)
		if doc.ID != "" {
			idStr := C.CString(doc.ID)
			defer C.free(unsafe.Pointer(idStr))
			C.zvec_doc_set_pk(cDoc, idStr)
		}

		// Add fields
		for name, value := range doc.Fields {
			if err := addFieldToDoc(cDoc, name, value); err != nil {
				return fmt.Errorf("failed to add field %s to doc %d: %v", name, i, err)
			}
		}
	}

	// Convert to C array
	cDocsArray := (**C.zvec_doc_t)(unsafe.Pointer(&cDocs[0]))

	// Call C API
	var successCount, errorCount C.size_t
	code := C.zvec_collection_insert(coll.handle, cDocsArray, C.size_t(len(docs)), &successCount, &errorCount)
	if code != C.ZVEC_OK {
		return fmt.Errorf("insert failed: %s (success: %d, errors: %d)", GetLastError(), successCount, errorCount)
	}

	return nil
}

// addFieldToDoc adds a field to a C document.
func addFieldToDoc(cDoc *C.zvec_doc_t, name string, value interface{}) error {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	switch v := value.(type) {
	case int32:
		return zvecErrorToGo(C.zvec_doc_add_field_by_value(cDoc, nameStr, C.ZVEC_DATA_TYPE_INT32, unsafe.Pointer(&v), 0))
	case int64:
		val := C.int64_t(v)
		return zvecErrorToGo(C.zvec_doc_add_field_by_value(cDoc, nameStr, C.ZVEC_DATA_TYPE_INT64, unsafe.Pointer(&val), 0))
	case float32:
		val := C.float(v)
		return zvecErrorToGo(C.zvec_doc_add_field_by_value(cDoc, nameStr, C.ZVEC_DATA_TYPE_FLOAT, unsafe.Pointer(&val), 0))
	case float64:
		val := C.double(v)
		return zvecErrorToGo(C.zvec_doc_add_field_by_value(cDoc, nameStr, C.ZVEC_DATA_TYPE_DOUBLE, unsafe.Pointer(&val), 0))
	case string:
		val := C.CString(v)
		defer C.free(unsafe.Pointer(val))
		return zvecErrorToGo(C.zvec_doc_add_field_by_value(cDoc, nameStr, C.ZVEC_DATA_TYPE_STRING, unsafe.Pointer(val), C.size_t(len(v))))
	case bool:
		val := C.bool(v)
		return zvecErrorToGo(C.zvec_doc_add_field_by_value(cDoc, nameStr, C.ZVEC_DATA_TYPE_BOOL, unsafe.Pointer(&val), 0))
	case []float32:
		// Vector field
		return fmt.Errorf("vector field not yet supported - use proper vector API")
	default:
		return fmt.Errorf("unsupported field type: %T", value)
	}
}

// UpsertDocuments upserts documents into a collection.
func UpsertDocuments(coll *Collection, docs []Document) error {
	// TODO: Implement using zvec_collection_upsert C API
	return fmt.Errorf("not implemented yet")
}

// DeleteDocuments deletes documents by ID.
func DeleteDocuments(coll *Collection, ids []string) error {
	// TODO: Implement using zvec_collection_delete C API
	return fmt.Errorf("not implemented yet")
}

// DeleteByFilter deletes documents matching a filter.
func DeleteByFilter(coll *Collection, filter string) error {
	// TODO: Implement using zvec_collection_delete_by_filter C API
	return fmt.Errorf("not implemented yet")
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID     string
	Score  float32
	Fields map[string]interface{}
}

// Search searches for similar vectors in a collection.
// This is a Go wrapper around the C API zvec_collection_query().
func Search(coll *Collection, query []float32, fieldName string, limit int) ([]SearchResult, error) {
	// TODO: Implement using zvec_collection_query C API
	// This requires:
	// 1. Creating query params (HNSW/IVF/Flat)
	// 2. Setting up the query
	// 3. Calling zvec_collection_query()
	// 4. Parsing results
	return nil, fmt.Errorf("not implemented yet")
}

// Query performs a filtered query.
func Query(coll *Collection, filter string, limit int) ([]Document, error) {
	// TODO: Implement using zvec_collection_query C API with filter
	return nil, fmt.Errorf("not implemented yet")
}

// MultiQuery performs multiple queries at once.
func MultiQuery(coll *Collection, queries []string, limit int) ([][]Document, error) {
	// TODO: Implement using zvec_collection_multi_query C API
	return nil, fmt.Errorf("not implemented yet")
}

// Fetch retrieves documents by ID.
func Fetch(coll *Collection, ids []string) ([]Document, error) {
	// TODO: Implement using zvec_collection_fetch C API
	return nil, fmt.Errorf("not implemented yet")
}

// GetStats returns collection statistics.
func GetStats(coll *Collection) (map[string]interface{}, error) {
	// TODO: Implement using zvec_collection_get_stats C API
	return nil, fmt.Errorf("not implemented yet")
}

// zvecErrorToGo converts a zvec error code to a Go error.
func zvecErrorToGo(code C.zvec_error_code_t) error {
	if code == C.ZVEC_OK {
		return nil
	}
	return fmt.Errorf("zvec error: %s", GetLastError())
}
