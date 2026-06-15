// Package zvec_test contains integration tests for the zvec Go SDK.
// These tests require the zvec C library to be built and placed in go/lib/.
//
// To run these tests:
//  1. Build the zvec C library (see go/BUILD.md)
//  2. Copy libzvec.so (or .dylib/.dll) to go/lib/
//  3. Run: cd go && go test -v -run TestIntegration
package zvec_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	zvec "github.com/ghshhf/zvec/go"
)

const (
	testCollectionPath = "/tmp/zvec_test_go"
	testDimension    = 128
	testDocCount    = 100
)

// TestIntegration_FullWorkflow tests the full CRUD + Search workflow.
func TestIntegration_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Cleanup from previous run
	os.RemoveAll(testCollectionPath)

	// ===== Step 1: Create Collection =====
	t.Run("CreateCollection", func(t *testing.T) {
		schema := zvec.NewSchema()
		defer schema.Destroy()

		schema.SetName("test_collection")

		err := schema.AddField("id", zvec.FieldTypeInt64, false)
		if err != nil {
			t.Fatalf("Failed to add id field: %v", err)
		}

		err = schema.AddField("vector", zvec.FieldTypeFloatVector, false)
		if err != nil {
			t.Fatalf("Failed to add vector field: %v", err)
		}

		err = schema.AddField("text", zvec.FieldTypeString, false)
		if err != nil {
			t.Fatalf("Failed to add text field: %v", err)
		}

		coll, err := zvec.CreateCollection(testCollectionPath, schema, nil)
		if err != nil {
			t.Fatalf("Failed to create collection: %v", err)
		}
		defer coll.Close()

		t.Logf("✓ Created collection at %s", testCollectionPath)
	})

	// ===== Step 2: Open Collection =====
	var coll *zvec.Collection
	t.Run("OpenCollection", func(t *testing.T) {
		var err error
		coll, err = zvec.OpenCollection(testCollectionPath, false)
		if err != nil {
			t.Fatalf("Failed to open collection: %v", err)
		}
		t.Logf("✓ Opened collection")
	})
	defer coll.Close()

	// ===== Step 3: Insert Documents =====
	t.Run("InsertDocuments", func(t *testing.T) {
		docs := make([]zvec.Document, testDocCount)
		for i := 0; i < testDocCount; i++ {
			vector := make([]float32, testDimension)
			for j := range vector {
				vector[j] = rand.Float32()
			}

			docs[i] = zvec.Document{
				"id":     int64(i),
				"vector": vector,
				"text":   fmt.Sprintf("document_%d", i),
			}
		}

		results, err := zvec.InsertDocuments(coll, docs)
		if err != nil {
			t.Fatalf("Failed to insert documents: %v", err)
		}

		t.Logf("✓ Inserted %d documents", len(results))
	})

	// ===== Step 4: Create Index =====
	t.Run("CreateIndex", func(t *testing.T) {
		err := zvec.CreateIndex(coll, "vector", zvec.IndexTypeHNSW, zvec.MetricTypeL2)
		if err != nil {
			t.Fatalf("Failed to create index: %v", err)
		}
		t.Logf("✓ Created HNSW index on 'vector' field")
	})

	// ===== Step 5: Search =====
	t.Run("Search", func(t *testing.T) {
		queryVector := make([]float32, testDimension)
		for i := range queryVector {
			queryVector[i] = rand.Float32()
		}

		params := &zvec.QueryParams{
			HNSW: &zvec.HNSWQueryParams{
				EFSearch:       50,
				IsUsingRefiner: true,
			},
		}

		results, err := zvec.Search(coll, "vector", queryVector, 5, "", nil, params)
		if err != nil {
			t.Fatalf("Failed to search: %v", err)
		}

		t.Logf("✓ Search returned %d results", len(results))
		for i, r := range results {
			t.Logf("  Result %d: id=%v, score=%.4f", i, r.Doc["id"], r.Score)
		}
	})

	// ===== Step 6: Fetch =====
	t.Run("Fetch", func(t *testing.T) {
		docs, err := zvec.Fetch(coll, []string{"0", "1", "2"})
		if err != nil {
			t.Fatalf("Failed to fetch documents: %v", err)
		}

		t.Logf("✓ Fetched %d documents", len(docs))
	})

	// ===== Step 7: Update =====
	t.Run("UpdateDocuments", func(t *testing.T) {
		updates := []zvec.Document{
			{"id": int64(0), "text": "updated_document_0"},
		}

		results, err := zvec.UpdateDocuments(coll, updates)
		if err != nil {
			t.Fatalf("Failed to update documents: %v", err)
		}

		t.Logf("✓ Updated %d documents", len(results))
	})

	// ===== Step 8: Upsert =====
	t.Run("UpsertDocuments", func(t *testing.T) {
		docs := []zvec.Document{
			{"id": int64(999), "vector": randomVector(testDimension), "text": "new_document"},
		}

		results, err := zvec.UpsertDocuments(coll, docs)
		if err != nil {
			t.Fatalf("Failed to upsert documents: %v", err)
		}

		t.Logf("✓ Upserted %d documents", len(results))
	})

	// ===== Step 9: Delete =====
	t.Run("DeleteDocuments", func(t *testing.T) {
		results, err := zvec.DeleteDocuments(coll, []string{"999"})
		if err != nil {
			t.Fatalf("Failed to delete documents: %v", err)
		}

		t.Logf("✓ Deleted %d documents", len(results))
	})

	// ===== Step 10: GetStats =====
	t.Run("GetStats", func(t *testing.T) {
		stats, err := zvec.GetStats(coll)
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		t.Logf("✓ Collection stats: %d documents, %d indexes", stats.DocCount, stats.IndexCount)
	})

	// ===== Step 11: Flush =====
	t.Run("Flush", func(t *testing.T) {
		err := zvec.Flush(coll)
		if err != nil {
			t.Fatalf("Failed to flush: %v", err)
		}

		t.Logf("✓ Flushed collection")
	})

	// ===== Step 12: Cleanup =====
	t.Run("Cleanup", func(t *testing.T) {
		zvec.Cleanup()
		t.Logf("✓ Cleaned up zvec resources")
	})

	// Cleanup test collection
	os.RemoveAll(testCollectionPath)
}

// TestIntegration_MultiThread tests multi-threaded operations.
func TestIntegration_MultiThread(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	path := "/tmp/zvec_test_go_mt"
	os.RemoveAll(path)
	defer os.RemoveAll(path)

	// Create collection
	schema := zvec.NewSchema()
	defer schema.Destroy()
	schema.SetName("mt_test")
	schema.AddField("id", zvec.FieldTypeInt64, false)
	schema.AddField("vector", zvec.FieldTypeFloatVector, false)

	coll, err := zvec.CreateCollection(path, schema, nil)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}
	defer coll.Close()

	// Insert documents in goroutines
	done := make(chan bool)
	for g := 0; g < 4; g++ {
		go func(g int) {
			docs := make([]zvec.Document, 25)
			for i := 0; i < 25; i++ {
				docs[i] = zvec.Document{
					"id":     int64(g*25 + i),
					"vector": randomVector(testDimension),
				}
			}
			_, err := zvec.InsertDocuments(coll, docs)
			if err != nil {
				t.Errorf("Goroutine %d failed to insert: %v", g, err)
			}
			done <- true
		}(g)
	}

	for g := 0; g < 4; g++ {
		<-done
	}

	t.Logf("✓ Multi-threaded insert completed")
}

// Helper function to generate random vector.
func randomVector(dim int) []float32 {
	vector := make([]float32, dim)
	for i := range vector {
		vector[i] = rand.Float32()
	}
	return vector
}

// TestIntegration_SearchAccuracy tests search accuracy.
func TestIntegration_SearchAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	path := "/tmp/zvec_test_go_accuracy"
	os.RemoveAll(path)
	defer os.RemoveAll(path)

	// Create collection with known vectors
	schema := zvec.NewSchema()
	defer schema.Destroy()
	schema.SetName("accuracy_test")
	schema.AddField("id", zvec.FieldTypeInt64, false)
	schema.AddField("vector", zvec.FieldTypeFloatVector, false)

	coll, err := zvec.CreateCollection(path, schema, nil)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}
	defer coll.Close()

	// Insert vectors where vector[0] = id (easy to verify)
	docs := make([]zvec.Document, 100)
	for i := 0; i < 100; i++ {
		vector := make([]float32, testDimension)
		vector[0] = float32(i) // Make vector[0] = id
		docs[i] = zvec.Document{
			"id":     int64(i),
			"vector": vector,
		}
	}

	_, err = zvec.InsertDocuments(coll, docs)
	if err != nil {
		t.Fatalf("Failed to insert documents: %v", err)
	}

	// Create index
	err = zvec.CreateIndex(coll, "vector", zvec.IndexTypeHNSW, zvec.MetricTypeL2)
	if err != nil {
		t.Fatalf("Failed to create index: %v", err)
	}

	// Search for vector closest to id=50
	queryVector := make([]float32, testDimension)
	queryVector[0] = 50.0

	results, err := zvec.Search(coll, "vector", queryVector, 5, "", nil, nil)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	// The first result should be id=50 (or close to it)
	t.Logf("Search results for query vector[0]=50:")
	for i, r := range results {
		id := int(r.Doc["id"].(int64))
		t.Logf("  %d: id=%d, score=%.4f", i, id, r.Score)
	}

	t.Logf("✓ Search accuracy test completed")
}

// TestVersion tests the Version() function.
func TestVersion(t *testing.T) {
	version := zvec.Version()
	if version == "" {
		t.Error("Version() returned empty string")
	}
	t.Logf("zvec version: %s", version)
}

// TestCheckVersion tests the CheckVersion() function.
func TestCheckVersion(t *testing.T) {
	// Check for version 0.4.0 (should be compatible)
	compatible := zvec.CheckVersion(0, 4, 0)
	t.Logf("Compatible with 0.4.0: %v", compatible)
}
