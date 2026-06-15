package zvec

/*
#include <zvec/c_api.h>
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
// TODO: Implement using C API - requires converting Go documents to C zvec_doc_t**
func InsertDocuments(coll *Collection, docs []Document) error {
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// UpsertDocuments upserts documents into a collection.
// TODO: Implement using C API
func UpsertDocuments(coll *Collection, docs []Document) error {
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// DeleteDocuments deletes documents by ID.
// TODO: Implement using C API
func DeleteDocuments(coll *Collection, ids []string) error {
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// DeleteByFilter deletes documents matching a filter.
// TODO: Implement using C API
func DeleteByFilter(coll *Collection, filter string) error {
	return fmt.Errorf("not implemented yet - requires C API binding")
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID     string
	Score  float32
	Fields map[string]interface{}
}

// Search searches for similar vectors in a collection.
// TODO: Implement using C API - requires setting up query params and calling search
func Search(coll *Collection, query []float32, fieldName string, limit int) ([]SearchResult, error) {
	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// Query performs a filtered query.
// TODO: Implement using C API
func Query(coll *Collection, filter string, limit int) ([]Document, error) {
	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// MultiQuery performs multiple queries at once.
// TODO: Implement using C API
func MultiQuery(coll *Collection, queries []string, limit int) ([][]Document, error) {
	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// Fetch retrieves documents by ID.
// TODO: Implement using C API
func Fetch(coll *Collection, ids []string) ([]Document, error) {
	return nil, fmt.Errorf("not implemented yet - requires C API binding")
}

// GetStats returns collection statistics.
func GetStats(coll *Collection) (map[string]interface{}, error) {
	// TODO: Implement using C API - call zvec_collection_get_stats
	return nil, fmt.Errorf("not implemented yet")
}
