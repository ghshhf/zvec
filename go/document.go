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
	FieldTypeInt32     FieldType = C.ZVEC_FIELD_TYPE_INT32
	FieldTypeInt64     FieldType = C.ZVEC_FIELD_TYPE_INT64
	FieldTypeFloat     FieldType = C.ZVEC_FIELD_TYPE_FLOAT
	FieldTypeDouble    FieldType = C.ZVEC_FIELD_TYPE_DOUBLE
	FieldTypeString    FieldType = C.ZVEC_FIELD_TYPE_STRING
	FieldTypeBinary    FieldType = C.ZVEC_FIELD_TYPE_BINARY
	FieldTypeBool      FieldType = C.ZVEC_FIELD_TYPE_BOOL
	FieldTypeJSON      FieldType = C.ZVEC_FIELD_TYPE_JSON
	FieldTypeVector    FieldType = C.ZVEC_FIELD_TYPE_VECTOR
	FieldTypeSparseVector FieldType = C.ZVEC_FIELD_TYPE_SPARSE_VECTOR
	FieldTypeDatetime  FieldType = C.ZVEC_FIELD_TYPE_DATETIME
)

// Document represents a document with fields.
type Document struct {
	ID     string
	Fields map[string]interface{}
}

// InsertDocuments inserts documents into a collection.
func InsertDocuments(coll *Collection, docs []Document) error {
	// TODO: Convert Go documents to C zvec_doc_t**
	// TODO: Call C API function
	return fmt.Errorf("not implemented yet")
}

// UpsertDocuments upserts documents into a collection.
func UpsertDocuments(coll *Collection, docs []Document) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// DeleteDocuments deletes documents by ID.
func DeleteDocuments(coll *Collection, ids []string) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// SearchResult represents a search result.
type SearchResult struct {
	ID    string
	Score float32
	Fields map[string]interface{}
}

// Search searches for similar vectors in a collection.
func Search(coll *Collection, query []float32, fieldName string, limit int) ([]SearchResult, error) {
	// TODO: Implement using C API
	return nil, fmt.Errorf("not implemented yet")
}

// Query performs a query with filters.
func Query(coll *Collection, filter string, limit int) ([]Document, error) {
	// TODO: Implement using C API
	return nil, fmt.Errorf("not implemented yet")
}
