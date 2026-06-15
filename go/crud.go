package zvec

/*
#include <zvec/c_api.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// DeleteDocuments deletes documents by their IDs.
// This is a Go wrapper around the C API zvec_collection_delete().
func DeleteDocuments(coll *Collection, ids []string) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}
	if len(ids) == 0 {
		return nil
	}

	// Convert Go string IDs to C strings
	cIDs := make([]*C.char, len(ids))
	defer func() {
		for _, id := range cIDs {
			if id != nil {
				C.free(unsafe.Pointer(id))
			}
		}
	}()

	for i, id := range ids {
		cIDs[i] = C.CString(id)
	}

	// Convert to C array
	cIDsArray := (**C.char)(unsafe.Pointer(&cIDs[0]))

	// Call C API
	var successCount, errorCount C.size_t
	code := C.zvec_collection_delete(coll.handle, cIDsArray, C.size_t(len(ids)), &successCount, &errorCount)
	if code != C.ZVEC_OK {
		return fmt.Errorf("delete failed: %s (success: %d, errors: %d)", GetLastError(), successCount, errorCount)
	}

	return nil
}

// DeleteByFilter deletes documents matching a filter expression.
// This is a Go wrapper around the C API zvec_collection_delete_by_filter().
func DeleteByFilter(coll *Collection, filter string) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}

	filterStr := C.CString(filter)
	defer C.free(unsafe.Pointer(filterStr))

	code := C.zvec_collection_delete_by_filter(coll.handle, filterStr)
	if code != C.ZVEC_OK {
		return fmt.Errorf("delete by filter failed: %s", GetLastError())
	}

	return nil
}

// Fetch retrieves documents by their IDs.
// This is a Go wrapper around the C API zvec_collection_fetch().
func Fetch(coll *Collection, ids []string) ([]Document, error) {
	if coll == nil || coll.handle == nil {
		return nil, fmt.Errorf("collection is nil")
	}
	if len(ids) == 0 {
		return []Document{}, nil
	}

	// Convert Go string IDs to C strings
	cIDs := make([]*C.char, len(ids))
	defer func() {
		for _, id := range cIDs {
			if id != nil {
				C.free(unsafe.Pointer(id))
			}
		}
	}()

	for i, id := range ids {
		cIDs[i] = C.CString(id)
	}

	// Convert to C array
	cIDsArray := (**C.char)(unsafe.Pointer(&cIDs[0]))

	// Call C API
	var cDocs **C.zvec_doc_t
	var docCount C.size_t
	code := C.zvec_collection_fetch(coll.handle, cIDsArray, C.size_t(len(ids)), &cDocs, &docCount)
	if code != C.ZVEC_OK {
		return nil, fmt.Errorf("fetch failed: %s", GetLastError())
	}

	// Convert C documents to Go documents
	docs := make([]Document, 0, int(docCount))
	for i := 0; i < int(docCount); i++ {
		cDoc := (*C.zvec_doc_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cDocs)) + uintptr(i)*unsafe.Sizeof(*cDocs)))
		doc, err := cDocToGo(cDoc)
		if err != nil {
			continue // Skip documents that fail to convert
		}
		docs = append(docs, *doc)
	}

	// Free C resources
	C.zvec_free_docs(cDocs, docCount)

	return docs, nil
}

// cDocToGo converts a C zvec_doc_t to a Go Document.
func cDocToGo(cDoc *C.zvec_doc_t) (*Document, error) {
	if cDoc == nil {
		return nil, fmt.Errorf("C document is nil")
	}

	doc := &Document{
		Fields: make(map[string]interface{}),
	}

	// Get PK (primary key = ID)
	pk := C.zvec_doc_get_pk(cDoc)
	if pk != nil {
		doc.ID = C.GoString(pk)
	}

	// TODO: Get fields from C document
	// This requires understanding the C API for iterating document fields

	return doc, nil
}

// UpsertDocuments upserts documents into a collection.
// This is a Go wrapper around the C API zvec_collection_upsert().
func UpsertDocuments(coll *Collection, docs []Document) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}
	if len(docs) == 0 {
		return nil
	}

	// Create C documents (similar to InsertDocuments)
	cDocs := make([]*C.zvec_doc_t, len(docs))
	defer func() {
		for _, doc := range cDocs {
			if doc != nil {
				C.zvec_doc_destroy(doc)
			}
		}
	}()

	for i, doc := range docs {
		cDoc := C.zvec_doc_create()
		if cDoc == nil {
			return fmt.Errorf("failed to create C document for doc %d", i)
		}
		cDocs[i] = cDoc

		// Set PK
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
	code := C.zvec_collection_upsert(coll.handle, cDocsArray, C.size_t(len(docs)), &successCount, &errorCount)
	if code != C.ZVEC_OK {
		return fmt.Errorf("upsert failed: %s (success: %d, errors: %d)", GetLastError(), successCount, errorCount)
	}

	return nil
}

// UpdateDocuments updates documents by their IDs.
// This is a Go wrapper around the C API zvec_collection_update().
func UpdateDocuments(coll *Collection, docs []Document) error {
	if coll == nil || coll.handle == nil {
		return fmt.Errorf("collection is nil")
	}
	if len(docs) == 0 {
		return nil
	}

	// Create C documents (similar to InsertDocuments)
	cDocs := make([]*C.zvec_doc_t, len(docs))
	defer func() {
		for _, doc := range cDocs {
			if doc != nil {
				C.zvec_doc_destroy(doc)
			}
		}
	}()

	for i, doc := range docs {
		cDoc := C.zvec_doc_create()
		if cDoc == nil {
			return fmt.Errorf("failed to create C document for doc %d", i)
		}
		cDocs[i] = cDoc

		// Set PK
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
	code := C.zvec_collection_update(coll.handle, cDocsArray, C.size_t(len(docs)), &successCount, &errorCount)
	if code != C.ZVEC_OK {
		return fmt.Errorf("update failed: %s (success: %d, errors: %d)", GetLastError(), successCount, errorCount)
	}

	return nil
}
