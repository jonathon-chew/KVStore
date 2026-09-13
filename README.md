# KVStore

An educational, in-memory key-value store built in Go. `KVStore` implements a hash table from scratch to explore the mechanics behind fast key lookup: deterministic hashing, bucket selection, collision handling, resizing, and a small CRUD-style API.

It is intentionally not a replacement for Go's `map`, Redis, or a database. The point is to make the underlying data structure visible and testable.

## Features

- String keys hashed with FNV-1a via Go's standard `hash/fnv` package
- Average constant-time lookup through bucket indexing
- Separate chaining for hash collisions
- Add, get, existence check, and remove operations
- Duplicate-key protection
- Automatic bucket growth at a `0.75` load factor, with entries rehashed into the larger table
- Values stored as `any`, allowing mixed value types within one running store
- A runnable example and Go test suite

## Quick Start

### Prerequisites

- Go `1.27.1` or a compatible Go toolchain

### Run the demonstration

```bash
cd Private/kvstore
go run ./cmd/KVStore
```

The demonstration adds `Example`, reads it back, removes it, then confirms it is no longer present:

```text
Example exists: true
Example was found with the value: Hello World
Example exists: false
```

### Run the tests

```bash
go test ./...
```

## Core Operations

The `HashTable` implementation currently lives in `cmd/KVStore/main.go` and exposes these methods within the demonstration package:

| Method | Behaviour |
| --- | --- |
| `Add(key, value)` | Stores a new key/value pair. Returns an error when the key already exists. |
| `Get(key)` | Returns the stored value and `true`, or `nil` and `false` when the key is absent. |
| `Remove(key)` | Removes an existing key. Returns an error when the key is absent. |
| `exists(key)` | Internal helper that finds the target bucket and reports whether the key exists. |

## How It Works

```text
key
 |
 v
FNV-1a hash
 |
 v
hash % bucket count
 |
 v
bucket [Entry, Entry, ...]
```

Each bucket is a slice of `Entry` values. When two keys resolve to the same bucket, both remain in that slice and are compared by key during lookup. This is separate chaining.

The table begins with a fixed number of buckets. Once `size / bucketCount` reaches `0.75`, the table doubles its bucket count and rehashes every entry because the bucket index depends on the number of buckets.

## Project Structure

```text
cmd/KVStore/
  main.go          hash table implementation and runnable example
  main_test.go     insertion and duplicate-key test coverage
internal/cli/
  cli.go           reserved for a future command-line interface
scripts/           development and CI helper scripts
```

## Current Status

The project is a learning implementation and runs entirely in memory. Data is not persisted, the hash table is not yet packaged for import by other Go programs, and the current executable uses a fixed example rather than accepting user input.

The existing tests cover successful insertion and duplicate-key rejection. Lookup, deletion, collision, and resize behaviour remain good candidates for assertion-based tests.

## Roadmap

- Extract the hash table into an importable package
- Add a constructor so callers do not initialise buckets directly
- Expand test coverage for lookup, removal, collisions, and resizing
- Add an interactive CLI in `internal/cli`
- Consider persistence and a small network API only after the in-memory API is stable

## License

See [LICENSE](./LICENSE) for details.
