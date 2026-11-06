package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/jonathon-chew/KVStore/internal/commands"
	"github.com/jonathon-chew/KVStore/internal/kvstore"
)

func handleConnection(conn net.Conn, kv *kvstore.HashTable) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		command, err := commands.ReadCommand(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Protocol error:", err)
			break
		}
		fmt.Printf("Command: /%s/ Args: /%v/\n", command.Name, command.Args)
		rv, parseErr := commands.ParseCommand(command.Name, command.Args, kv)
		if parseErr != nil {
			conn.Write([]byte("-" + parseErr.Error() + "\r\n"))
			continue
		}

		if rv != nil {
			conn.Write([]byte(fmt.Sprintf("+%v\r\n", rv)))
		} else {
			conn.Write([]byte("+Ok\r\n"))
		}

	}
}

func main() {

	fmt.Println("Up and running...")

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
