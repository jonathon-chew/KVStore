# KVStore

KVStore is an in-memory key-value server written in Go. It accepts line-delimited TCP commands and stores values in a custom hash table.

## What it does

- Listens on TCP port `6379` on all network interfaces
- Accepts one command per newline from each connected client
- Supports `SET`, `GET`, and `DEL` commands
- Stores values in memory only; restarting the server clears all data
- Uses FNV-1a hashing and separate chaining for hash-table collisions
- Rejects attempts to `SET` an existing key
- Protects hash-table operations with a mutex so multiple connections can access the store
- Doubles the bucket count and rehashes entries once the tracked load factor reaches `0.75`

## Requirements

- Go `1.27.1`, as declared in [go.mod](./go.mod)
- A TCP client such as `nc` (netcat) for manual interaction

## Run the server

From the repository root:

```bash
go run ./cmd/KVStore
```

In another terminal, connect locally:

```bash
nc localhost 6379
```

## Command protocol

Send a single command followed by a newline. Commands are uppercase and use spaces to separate arguments.

| Command | Description | Success response | Error response |
| --- | --- | --- | --- |
| `SET <key> <value>` | Adds a new key with a string value. | `Successfully added` | `<key> already exists` |
| `GET <key>` | Retrieves a stored value. | The stored string value | `[ERROR]: No key for: <key>, false` |
| `DEL <key>` | Removes a stored key. | `Successfully added` | `<key> does not exist` |

For example:

```text
SET greeting Hello
Successfully added
GET greeting
Hello
DEL greeting
Successfully added
GET greeting
[ERROR]: No key for: greeting, false
```

`SET` values are single tokens, so they cannot contain spaces.

## Implementation notes

The store lives in [internal/kvstore/KV.go](./internal/kvstore/KV.go). Each hash-table bucket is a slice of entries; entries with the same bucket index are compared by key. The server in [cmd/KVStore/main.go](./cmd/KVStore/main.go) starts a goroutine for every client connection, while [internal/commands/parse_command.go](./internal/commands/parse_command.go) translates protocol commands into store operations.

## Test

```bash
go test ./...
```

## License

KVStore is released under the [MIT License](./LICENSE).
