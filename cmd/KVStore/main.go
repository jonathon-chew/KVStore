package main

import (
	"os"

	"github.com/jonathon-chew/KVStore/pkg/server"
)

func main() {
	os.Exit(server.Server(""))
}
