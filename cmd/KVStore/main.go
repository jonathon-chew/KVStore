package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/jonathon-chew/KVStore/internal/commands"
	"github.com/jonathon-chew/KVStore/internal/kvstore"
)

func handleConnection(conn net.Conn, kv *kvstore.HashTable) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		request := scanner.Text()
		fmt.Printf("Recieved: %s\n", request)

		response, err := commands.ParseCommand(request, kv)
		if err != nil {
			fmt.Println(err.Error())
			conn.Write([]byte(err.Error() + "\n"))
		}

		if r, ok := response.(string); ok {
			conn.Write([]byte(r + "\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read error:", err)
	}
}

func main() {
	var kv kvstore.HashTable
	kv.Buckets = make([][]kvstore.Entry, 10)

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[ERROR]: %s\n", err.Error())
		}

		go handleConnection(conn, &kv)
	}
}
