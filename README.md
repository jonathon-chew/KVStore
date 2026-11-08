# KVStore

KVStore is an in-memory key-value server written in Go. It accepts Redis-style RESP commands over TCP and stores values in a custom hash table.

## What it does

- Listens on TCP port `6379` on all network interfaces
- Accepts RESP arrays containing bulk strings from each connected client
- Supports `SET`, `GET`, `DEL`, `GETDEL`, `INCR`, and `STRLEN`
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

Send commands using the Redis Serialization Protocol (RESP). The first array item is the uppercase command and the remaining items are its arguments.

| Command | Description | Success response | Error response |
| --- | --- | --- | --- |
| `SET <key> <value>` | Adds a new key with a string value. | `Successfully added` | `<key> already exists` |
| `GET <key>` | Retrieves a stored value. | The stored string value | `[ERROR]: No key for: <key>, false` |
| `DEL <key>` | Removes a stored key. | `Successfully added` | `<key> does not exist` |

For example:

```text
*3\r\n$3\r\nSET\r\n$8\r\ngreeting\r\n$5\r\nHello\r\n
*2\r\n$3\r\nGET\r\n$8\r\ngreeting\r\n
Hello
```

`SET` values are single tokens, so they cannot contain spaces.

## Implementation notes

The store lives in [pkg/kvstore/KV.go](./pkg/kvstore/KV.go). Each hash-table bucket is a slice of entries; entries with the same bucket index are compared by key. The server in [pkg/server/server.go](./pkg/server/server.go) handles client connections, while [pkg/commands/parse_command.go](./pkg/commands/parse_command.go) reads RESP and translates commands into store operations.

## Test

```bash
go test ./...
```

## Roadmap
[x] HashMap
[x] Resize
[x] Thread safety
[x] TCP
[x] Protocol separation
[x] SET
[x] GET
[x] DEL
[x] GETDEL
[x] INCR
[x] STRLEN

[x] RESP
[ ] CLI
[ ] Pipelining
[ ] TTL + expiry heap
[ ] AOF/WAL
[ ] Snapshots
[ ] Metrics
[ ] Benchmarks/load testing

## License

KVStore is released under the [MIT License](./LICENSE).
