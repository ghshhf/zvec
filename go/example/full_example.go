package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	zvec "github.com/ghshhf/zvec/go"
)

func main() {
	fmt.Println("🚀 zvec Go SDK Example")
	fmt.Println("===============================")

	// 1. Check version
	version := zvec.Version()
	fmt.Printf("✅ zvec version: %s\n", version)

	// 2. Set log level
	zvec.SetLogLevel(zvec.LogLevelInfo)
	fmt.Println("✅ Log level set to INFO")

	// 3. Create a schema
	fmt.Println("\n📋 Creating schema...")
	schema, err := zvec.CreateSchema("my_collection")
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Add fields
	fields := map[string]zvec.FieldType{
		"id":        zvec.FieldTypeInt64,
		"name":      zvec.FieldTypeString,
		"embedding": zvec.FieldTypeVector,
		"metadata":  zvec.FieldTypeJSON,
		"timestamp": zvec.FieldTypeDatetime,
	}

	for name, fieldType := range fields {
		if err := schema.AddField(name, fieldType); err != nil {
			log.Fatalf("Failed to add field %s: %v", name, err)
		}
		fmt.Printf("  ✅ Added field: %s (type: %d)\n", name, fieldType)
	}

	// 4. Create options
	fmt.Println("\n⚙️  Creating options...")
	opts, err := zvec.CreateOptions()
	if err != nil {
		log.Fatalf("Failed to create options: %v", err)
	}

	if err := opts.SetEnableMMap(true); err != nil {
		log.Fatalf("Failed to set mmap: %v", err)
	}
	fmt.Println("  ✅ MMap enabled")

	// 5. Create collection
	dbPath := filepath.Join(".", "example_db")
	defer os.RemoveAll(dbPath)

	fmt.Printf("\n💾 Creating collection at: %s\n", dbPath)
	coll, err := zvec.CreateCollection(dbPath, schema, opts)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	defer coll.Close()
	fmt.Println("  ✅ Collection created successfully!")

	// 6. Get schema (from collection)
	fmt.Println("\n📖 Getting collection schema...")
	schema2, err := coll.GetSchema()
	if err != nil {
		log.Fatalf("Failed to get schema: %s", err)
	}
	fmt.Println("  ✅ Schema retrieved successfully")

	// Check if field exists
	if zvec.HasField(schema2, "embedding") {
		fmt.Println("  ✅ Field 'embedding' exists")
	}

	// 7. Validate schema
	fmt.Println("\n✔️  Validating schema...")
	if err := zvec.ValidateSchema(schema2); err != nil {
		fmt.Printf("  ⚠️  Schema validation warning: %v\n", err)
	} else {
		fmt.Println("  ✅ Schema is valid")
	}

	// 8. TODO: Insert documents
	fmt.Println("\n📝 TODO: Insert documents...")
	fmt.Println("  ⚡ This requires C API binding implementation")
	fmt.Println("  ⚡ Once implemented, you can:")
	fmt.Println("     doc := zvec.Document{")
	fmt.Println(`       ID: "doc1",`)
	fmt.Println("       Fields: map[string]interface{}{")
	fmt.Println(`         "name": "Example document",`)
	fmt.Println("         ...}")
	fmt.Println("     }")
	fmt.Println("     zvec.InsertDocuments(coll, []zvec.Document{doc})")

	// 9. TODO: Create index
	fmt.Println("\n🔍 TODO: Create vector index...")
	fmt.Println("  ⚡ This requires C API binding implementation")
	fmt.Println("  ⚡ Once implemented, you can:")
	fmt.Println("     zvec.CreateIndex(coll, \"embedding\", zvec.IndexTypeHNSW, zvec.MetricTypeCosine)")

	// 10. TODO: Search
	fmt.Println("\n🔎 TODO: Search similar vectors...")
	fmt.Println("  ⚡ This requires C API binding implementation")
	fmt.Println("  ⚡ Once implemented, you can:")
	fmt.Println("     results, err := zvec.Search(coll, queryVector, \"embedding\", 10)")

	fmt.Println("\n===============================")
	fmt.Println("🎉 Example completed successfully!")
	fmt.Println("🚧 Go SDK is still under development.")
	fmt.Println("📢 Once C API bindings are implemented, full functionality will be available.")
}
