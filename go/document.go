package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Document represents a document with vector and scalar fields.
type Document struct {
	ID     string
	Fields map[string]interface{}
}

// InsertDocuments inserts documents into a collection.
func InsertDocuments(collection *Collection, docs []Document) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// UpsertDocuments upserts documents into a collection.
func UpsertDocuments(collection *Collection, docs []Document) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// DeleteDocuments deletes documents by ID.
func DeleteDocuments(collection *Collection, ids []string) error {
	// TODO: Implement using C API
	return fmt.Errorf("not implemented yet")
}

// SearchResult represents a search result.
type SearchResult struct {
	ID       string
	Score    float32
	Fields   map[string]interface{}
}

// Search searches for similar vectors.
func Search(collection *Collection, query []float32, limit int) ([]SearchResult, error) {
	// TODO: Implement using C API
	return nil, fmt.Errorf("not implemented yet")
}
