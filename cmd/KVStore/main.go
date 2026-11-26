package main

import (
	"os"

	"github.com/jonathon-chew/KVStore/pkg/server"
)

// main starts the KVStore server with its default configuration.
func main() {
	os.Exit(server.Server(""))
}
