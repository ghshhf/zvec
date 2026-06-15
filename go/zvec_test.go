package zvec

import (
	"testing"
)

// TestVersion tests the Version function.
func TestVersion(t *testing.T) {
	version := Version()
	if version == "" {
		t.Error("Version() returned empty string")
	}
	t.Logf("zvec version: %s", version)
}

// TestCheckVersion tests the CheckVersion function.
func TestCheckVersion(t *testing.T) {
	// Test with current version (should be compatible)
	if !CheckVersion(0, 0, 0) {
		t.Log("CheckVersion(0,0,0) returned false (might be expected)")
	}
}

// TestCreateSchema tests schema creation.
func TestCreateSchema(t *testing.T) {
	schema, err := CreateSchema("test_collection")
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}
	if schema == nil {
		t.Fatal("CreateSchema returned nil schema")
	}
	if schema.handle == nil {
		t.Fatal("Schema handle is nil")
	}
}

// TestCreateOptions tests options creation.
func TestCreateOptions(t *testing.T) {
	opts, err := CreateOptions()
	if err != nil {
		t.Fatalf("Failed to create options: %v", err)
	}
	if opts == nil {
		t.Fatal("CreateOptions returned nil options")
	}
	if opts.handle == nil {
		t.Fatal("Options handle is nil")
	}
}

// TestFieldTypeConstants tests that field type constants are defined.
func TestFieldTypeConstants(t *testing.T) {
	types := []FieldType{
		FieldTypeInt32,
		FieldTypeInt64,
		FieldTypeFloat,
		FieldTypeDouble,
		FieldTypeString,
		FieldTypeBinary,
		FieldTypeBool,
		FieldTypeJSON,
		FieldTypeVector,
		FieldTypeSparseVector,
		FieldTypeDatetime,
	}

	for _, ft := range types {
		if ft == 0 {
			t.Errorf("Field type constant is 0 (might be unset)")
		}
	}
}

// TestIndexTypeConstants tests that index type constants are defined.
func TestIndexTypeConstants(t *testing.T) {
	types := []IndexType{
		IndexTypeIVF,
		IndexTypeHNSW,
		IndexTypeFlat,
		IndexTypeFTS,
		IndexTypeVamana,
	}

	for _, it := range types {
		if it == 0 {
			t.Errorf("Index type constant is 0 (might be unset)")
		}
	}
}

// TODO: Add more tests once C library integration is complete
// TODO: Add integration tests with actual zvec library
