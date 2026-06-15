# zvec Go SDK

Go bindings for [zvec](https://github.com/alibaba/zvec) - a lightweight, lightning-fast vector database.

## Status

✅ **Core functionality implemented** - This Go SDK provides complete bindings for the zvec C API.

- ✅ Collection create/open/close/destroy
- ✅ Schema management (create, add field, validate)
- ✅ Options management (mmap, read-only, buffer size)
- ✅ Document operations (insert, upsert, update, delete, fetch)
- ✅ Index management (create, drop, list)
- ✅ Vector search (with HNSW/IVF/Flat query parameters)
- ✅ Batch operations (batch insert, batch search)
- ✅ utility functions (flush, stats, compact, validate)
- ✅ Complete type definitions (FieldType, IndexType, MetricType)
- ✅ Examples and integration tests

## Requirements

- Go 1.21+
- zvec C library built and installed (`libzvec.so` / `zvec.dll` / `libzvec.dylib`)
- CGO enabled

## Building

See [BUILD.md](./BUILD.md) for detailed build instructions.

Quick start:
```bash
# 1. Build zvec C library
cd /path/to/zvec
mkdir -p build && cd build
cmake .. && make -j$(nproc)

# 2. Copy library to Go SDK
cp build/libzvec.so go/lib/              # Linux
cp build/libzvec.dylib go/lib/          # macOS
cp build/zvec.dll go/lib/               # Windows

# 3. Copy header
cp src/include/zvec/c_api.h go/include/zvec/

# 4. Run tests
cd go
go test -v
```

## Installation

```bash
go get github.com/ghshhf/zvec/go
```

## Usage

### Create a Collection

```go
package main

import (
    "fmt"
    "log"

    zvec "github.com/ghshhf/zvec/go"
)

func main() {
    // Create schema
    schema := zvec.NewSchema()
    defer schema.Destroy()

    schema.SetName("my_collection")
    schema.AddField("id", zvec.FieldTypeInt64, false)
    schema.AddField("embedding", zvec.FieldTypeFloatVector, false)
    schema.AddField("text", zvec.FieldTypeString, false)

    // Create collection
    coll, err := zvec.CreateCollection("./my_db", schema, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer coll.Close()

    fmt.Println("Collection created successfully!")
}
```

### Insert Documents

```go
// Create documents
docs := []zvec.Document{
    {
        "id":       int64(1),
        "embedding": []float32{0.1, 0.2, 0.3, ..., 0.128}, // 128-dim vector
        "text":     "hello world",
    },
    // ... more documents
}

// Insert
results, err := zvec.InsertDocuments(coll, docs)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Inserted %d documents\n", len(results))
```

### Create Index and Search

```go
// Create HNSW index
err = zvec.CreateIndex(coll, "embedding", zvec.IndexTypeHNSW, zvec.MetricTypeL2)
if err != nil {
    log.Fatal(err)
}

// Search
queryVector := []float32{0.1, 0.2, 0.3, ..., 0.128}
params := &zvec.QueryParams{
    HNSW: &zvec.HNSWQueryParams{
        EFSearch:      50,
        IsUsingRefiner: true,
    },
}

results, err := zvec.Search(coll, "embedding", queryVector, 5, "", nil, params)
if err != nil {
    log.Fatal(err)
}

for i, r := range results {
    fmt.Printf("Result %d: id=%v, score=%.4f\n", i, r.Doc["id"], r.Score)
}
```

### Fetch and Update

```go
// Fetch documents by ID
docs, err := zvec.Fetch(coll, []string{"1", "2", "3"})
if err != nil {
    log.Fatal(err)
}

// Update
updates := []zvec.Document{
    {"id": int64(1), "text": "updated text"},
}
_, err = zvec.UpdateDocuments(coll, updates)
if err != nil {
    log.Fatal(err)
}
```

### Delete and Cleanup

```go
// Delete by IDs
_, err = zvec.DeleteDocuments(coll, []string{"1"})
if err != nil {
    log.Fatal(err)
}

// Flush to disk
err = zvec.Flush(coll)
if err != nil {
    log.Fatal(err)
}

// Get stats
stats, err := zvec.GetStats(coll)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Documents: %d, Indexes: %d\n", stats.DocCount, stats.IndexCount)
```

## Examples

See `example/main.go` and `example/full_example.go` for complete working examples.

Run examples:
```bash
cd go/example
go run main.go
go run full_example.go
```

## Testing

Run unit tests:
```bash
cd go
go test -v
```

Run integration tests (requires C library):
```bash
cd go
go test -v -run TestIntegration
```

## TODO

- [ ] Complete score/distance in SearchResult (requires C API update)
- [ ] Add iterator support
- [ ] Add transaction support
- [ ] Performance optimization
- [ ] Windows CI/CD pipeline

## License

Apache License 2.0 (same as zvec)
