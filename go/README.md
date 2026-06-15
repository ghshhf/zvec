# zvec Go SDK

Go bindings for [zvec](https://github.com/alibaba/zvec) - a lightweight, lightning-fast vector database.

## Status

🚧 **Work in Progress** - This is an early-stage Go SDK implementation.

- ✅ Basic structure and types defined
- ✅ Core C API bindings framework
- ⚡ Collection create/open/close/destroy
- ⚡ Schema management (add field, set name)
- ⚡ Options management (mmap, read-only)
- ❗ Document operations (insert, update, delete, search) - TODO
- ❗ Index management (create, drop) - TODO
- ❗ Query and filter operations - TODO

## Requirements

- Go 1.21+
- zvec C library installed (`libzvec.so` / `zvec.dll` / `libzvec.dylib`)
- CGO enabled

## Installation

```bash
go get github.com/ghshhf/zvec/go
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "zvec"
)

func main() {
    // Create a schema
    schema, err := zvec.CreateSchema("my_collection")
    if err != nil {
        log.Fatal(err)
    }

    // Add fields
    schema.AddField("id", zvec.FieldTypeInt64)
    schema.AddField("embedding", zvec.FieldTypeVector)

    // Create collection
    coll, err := zvec.CreateCollection("./my_db", schema, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer coll.Close()

    fmt.Println("Collection created successfully!")
}
```

## Building

1. Build zvec C library first (see main zvec README)
2. Copy `libzvec` to `go/lib/`
3. Copy C headers to `go/include/`
4. Run `go build ./go/...`

## TODO

- [ ] Complete document operations (Insert, Update, Delete, Upsert)
- [ ] Complete index operations (CreateIndex, DropIndex)
- [ ] Complete search and query operations
- [ ] Add batch operations support
- [ ] Add iterator support
- [ ] Add proper memory management (free C strings/arrays)
- [ ] Add comprehensive tests
- [ ] Add examples
- [ ] Add documentation

## License

Apache License 2.0 (same as zvec)
