package main

import (
	"fmt"
	"log"
	"os"

	zvec "github.com/ghshhf/zvec/go"
)

func main() {
	// Create a schema
	schema, err := zvec.CreateSchema("example_collection")
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Add fields
	if err := schema.AddField("id", zvec.FieldTypeInt64); err != nil {
		log.Fatalf("Failed to add field: %v", err)
	}

	if err := schema.AddField("name", zvec.FieldTypeString); err != nil {
		log.Fatalf("Failed to add field: %v", err)
	}

	if err := schema.AddField("embedding", zvec.FieldTypeVector); err != nil {
		log.Fatalf("Failed to add field: %v", err)
	}

	// Create collection
	dbPath := "./example_db"
	defer os.RemoveAll(dbPath)

	coll, err := zvec.CreateCollection(dbPath, schema, nil)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	defer coll.Close()

	fmt.Println("✅ Collection created successfully!")

	// Get version
	version := zvec.Version()
	fmt.Printf("zvec version: %s\n", version)

	// TODO: Add more examples once document operations are implemented
	fmt.Println("🚧 Go SDK is still under development. More features coming soon!")
}
